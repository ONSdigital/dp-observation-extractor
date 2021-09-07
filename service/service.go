package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/initialise"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-reporter-client/reporter"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

func Run(ctx context.Context, config *config.Config, serviceList initialise.ExternalServiceList, signals chan os.Signal, errorChannel chan error, BuildTime string, GitCommit string, Version string) error {

	// S3 Session and clients (mapped by bucket name)
	sess, s3Clients, err := serviceList.GetS3Clients(config)
	if err != nil {
		return err
	}

	// Kafka Consumer
	kafkaConsumer, err := serviceList.GetConsumer(ctx, config)
	if err != nil {
		return err
	}

	// Kafka Observation Producer
	kafkaObservationProducer, err := serviceList.GetProducer(ctx, config.ObservationProducerTopic, initialise.Observation, config)
	if err != nil {
		return err
	}

	// Kafka Error Reporter
	kafkaErrorProducer, err := serviceList.GetProducer(ctx, config.ErrorProducerTopic, initialise.ErrorReporter, config)
	if err != nil {
		return err
	}

	observationWriter := observation.NewMessageWriter(kafkaObservationProducer)

	// Vault Client
	var vaultClient event.VaultClient
	if !config.EncryptionDisabled {
		vaultClient, err = vault.CreateClient(config.VaultToken, config.VaultAddr, 3)
		if err != nil {
			return err
		}
	}

	// Create healthcheck object with versionInfo
	hc, err := serviceList.GetHealthChecker(ctx, BuildTime, GitCommit, Version, config)
	if err != nil {
		return err
	}

	registerCheckers(ctx, hc, kafkaConsumer, kafkaObservationProducer, kafkaErrorProducer, vaultClient, s3Clients)
	if err != nil {
		return err
	}

	httpServer := startHealthCheck(ctx, hc, config.BindAddr, errorChannel)

	eventHandler := event.NewCSVHandler(sess, s3Clients, vaultClient, observationWriter, config.VaultPath)

	errorReporter, err := reporter.NewImportErrorReporter(kafkaErrorProducer, log.Namespace)
	if err != nil {
		return err
	}

	eventConsumer := event.NewConsumer()
	eventConsumer.Consume(ctx, kafkaConsumer, eventHandler, errorReporter)

	shutdownGracefully := func() error {

		ctx, cancel := context.WithTimeout(ctx, config.GracefulShutdownTimeout)
		anyError := false

		// Stop listening to Kafka consumer
		if serviceList.Consumer {
			if err = kafkaConsumer.StopListeningToConsumer(ctx); err != nil {
				anyError = true
				log.Error(ctx, "bad kafka consumer listen stop", err, log.Data{"topic": config.FileConsumerTopic})
			} else {
				log.Info(ctx, "kafka consumer listen stopped", log.Data{"topic": config.FileConsumerTopic})
			}
		}

		// Stop HTTP server
		if err = httpServer.Shutdown(ctx); err != nil {
			anyError = true
			log.Error(ctx, "bad http server stop", err)
		} else {
			log.Info(ctx, "http server stopped")
		}

		// Stop healthcheck
		hc.Stop()
		log.Info(ctx, "healthcheck stopped")

		// Close event consumer
		if err = eventConsumer.Close(ctx); err != nil {
			anyError = true
			log.Error(ctx, "bad event consumer stop", err)
		} else {
			log.Info(ctx, "event consumer stopped")
		}

		// Close Kafka consumer
		if serviceList.Consumer {
			if err = kafkaConsumer.Close(ctx); err != nil {
				anyError = true
				log.Error(ctx, "bad kafka consumer stop", err, log.Data{"topic": config.FileConsumerTopic})
			} else {
				log.Info(ctx, "kafka consumer stopped", log.Data{"topic": config.FileConsumerTopic})
			}
		}

		// Close Error Reporter Kafka producer
		if serviceList.ErrorReporterProducer {
			if err = kafkaErrorProducer.Close(ctx); err != nil {
				anyError = true
				log.Error(ctx, "bad kafka error reporter producer stop", err, log.Data{"topic": config.ErrorProducerTopic})
			} else {
				log.Info(ctx, "kafka error report producer stopped", log.Data{"topic": config.ErrorProducerTopic})
			}
		}

		// Close Observation Kafka producer
		if serviceList.ObservationProducer {
			if err = kafkaObservationProducer.Close(ctx); err != nil {
				anyError = true
				log.Error(ctx, "bad kafka observation producer stop", err, log.Data{"topic": config.ObservationProducerTopic})
			} else {
				log.Info(ctx, "kafka observation producer stopped", log.Data{"topic": config.ObservationProducerTopic})
			}
		}

		// cancel the timer in the shutdown context.
		cancel()

		// if any error happened during shutdown, log it and exit with err code
		if anyError {
			log.Warn(ctx, "graceful shutdown had errors")
			return errors.New("Failed to shutdown gracefully")
		}

		// if all dependencies shutted down successfully, log it and exit with success code
		log.Info(ctx, "graceful shutdown was successful")
		return nil
	}

	// Log non-fatal errors in separate go routines
	kafkaConsumer.Channels().LogErrors(ctx, "kafka consumer error")
	kafkaObservationProducer.Channels().LogErrors(ctx, "kafka observation producer error")
	kafkaErrorProducer.Channels().LogErrors(ctx, "kafka error producer error")
	go func() {
		for {
			select {
			case err := <-errorChannel:
				log.Error(ctx, "error channel", err)
			}
		}
	}()

	// When a signal is received, shutdown gracefully
	<-signals
	log.Error(ctx, "os signal received", errors.New("os signal received"))

	return shutdownGracefully()
}

// StartHealthCheck sets up the Handler, starts the healthcheck and the http server that serves health endpoint
func startHealthCheck(ctx context.Context, hc *healthcheck.HealthCheck, bindAddr string, errorChannel chan error) *server.Server {
	router := mux.NewRouter()
	router.Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	httpServer := server.New(bindAddr, router)
	httpServer.HandleOSSignals = false

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error(ctx, "http server error", err)
			hc.Stop()
			errorChannel <- err
		}
	}()
	return httpServer
}

// RegisterCheckers adds the checkers for the provided clients to the healthcheck object.
func registerCheckers(ctx context.Context, hc *healthcheck.HealthCheck,
	kafkaConsumer *kafka.ConsumerGroup,
	kafkaObservationProducer *kafka.Producer,
	kafkaErrorProducer *kafka.Producer,
	vaultClient event.VaultClient,
	s3Clients map[string]event.S3Client) (err error) {

	hasErrors := false

	if err = hc.AddCheck("Kafka Consumer", kafkaConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka consumer checker", err)
	}

	if err = hc.AddCheck("Kafka Observation Producer", kafkaObservationProducer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka observation producer checker", err)
	}

	if err = hc.AddCheck("Kafka Error Producer", kafkaErrorProducer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka error producer checker", err)
	}

	if vaultClient != nil {
		if err = hc.AddCheck("Vault", vaultClient.Checker); err != nil {
			hasErrors = true
			log.Error(ctx, "error adding check for vault checker", err)
		}
	}

	for bucketName, s3 := range s3Clients {
		if err := hc.AddCheck(fmt.Sprintf("S3 bucket %s", bucketName), s3.Checker); err != nil {
			hasErrors = true
			log.Error(ctx, "error adding check for s3 client", err)
		}
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
