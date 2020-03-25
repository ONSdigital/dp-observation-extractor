package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka"
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/initialise"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-reporter-client/reporter"
	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {

	log.Namespace = "dp-observation-extractor"
	ctx := context.Background()

	config, err := config.Get()
	if err != nil {
		log.Event(ctx, "error getting config", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	// sensitive fields are omitted from config.String().
	log.Event(ctx, "config on startup", log.INFO, log.Data{"config": config})

	// a channel used to signal a graceful exit is required.
	errorChannel := make(chan error)

	// Signal channel to know if SIGTERM is triggered
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// serviceList keeps track of what dependency services have been initialised
	serviceList := initialise.ExternalServiceList{}

	// S3 Session and clients (mapped by bucket name)
	sess, s3Clients, err := serviceList.GetS3Clients(config)
	checkForError(ctx, err)

	// Kafka Consumer
	kafkaConsumer, err := serviceList.GetConsumer(ctx, config)
	checkForError(ctx, err)

	// Kafka Observation Producer
	kafkaObservationProducer, err := serviceList.GetProducer(ctx, config.KafkaAddr, config.ObservationProducerTopic, initialise.Observation)
	checkForError(ctx, err)

	// Kafka Error Reporter
	kafkaErrorProducer, err := serviceList.GetProducer(ctx, config.KafkaAddr, config.ErrorProducerTopic, initialise.ErrorReporter)
	checkForError(ctx, err)

	observationWriter := observation.NewMessageWriter(kafkaObservationProducer)

	// Vault Client
	var vaultClient event.VaultClient
	if !config.EncryptionDisabled {
		vaultClient, err = vault.CreateClient(config.VaultToken, config.VaultAddr, 3)
		checkForError(ctx, err)
	}

	// Create healthcheck object with versionInfo
	hc, err := serviceList.GetHealthChecker(ctx, BuildTime, GitCommit, Version, config)
	checkForError(ctx, err)

	registerCheckers(ctx, hc, kafkaConsumer, kafkaObservationProducer, kafkaErrorProducer, vaultClient, s3Clients)
	checkForError(ctx, err)

	httpServer := startHealthCheck(ctx, hc, config.BindAddr)

	eventHandler := event.NewCSVHandler(sess, s3Clients, vaultClient, observationWriter, config.VaultPath)

	errorReporter, err := reporter.NewImportErrorReporter(kafkaErrorProducer, log.Namespace)
	checkForError(ctx, err)

	eventConsumer := event.NewConsumer()
	eventConsumer.Consume(ctx, kafkaConsumer, eventHandler, errorReporter)

	shutdownGracefully := func() {

		ctx, cancel := context.WithTimeout(ctx, config.GracefulShutdownTimeout)

		// Stop listening to Kafka consumer
		if serviceList.Consumer {
			if err = kafkaConsumer.StopListeningToConsumer(ctx); err != nil {
				log.Event(ctx, "bad kafka consumer listen stop", log.ERROR, log.Error(err), log.Data{"topic": config.FileConsumerTopic})
			} else {
				log.Event(ctx, "kafka consumer listen stopped", log.INFO, log.Data{"topic": config.FileConsumerTopic})
			}
		}

		// Stop HTTP server
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Event(ctx, "bad http server stop", log.ERROR, log.Error(err))
		} else {
			log.Event(ctx, "http server stopped", log.INFO)
		}

		// Stop healthcheck
		hc.Stop()

		// Close event consumer
		if err = eventConsumer.Close(ctx); err != nil {
			log.Event(ctx, "bad event consumer stop", log.ERROR, log.Error(err))
		} else {
			log.Event(ctx, "event consumer stopped", log.INFO)
		}

		// Close Kafka consumer
		if serviceList.Consumer {
			if err = kafkaConsumer.Close(ctx); err != nil {
				log.Event(ctx, "bad kafka consumer stop", log.ERROR, log.Error(err), log.Data{"topic": config.FileConsumerTopic})
			} else {
				log.Event(ctx, "kafka consumer stopped", log.INFO, log.Data{"topic": config.FileConsumerTopic})
			}
		}

		// Close Error Reporter Kafka producer
		if serviceList.ErrorReporterProducer {
			if err = kafkaErrorProducer.Close(ctx); err != nil {
				log.Event(ctx, "bad kafka error reporter producer stop", log.ERROR, log.Error(err), log.Data{"topic": config.ErrorProducerTopic})
			} else {
				log.Event(ctx, "kafka error report producer stopped", log.INFO, log.Data{"topic": config.ErrorProducerTopic})
			}
		}

		// Close Observation Kafka producer
		if serviceList.ObservationProducer {
			if err = kafkaObservationProducer.Close(ctx); err != nil {
				log.Event(ctx, "bad kafka observation producer stop", log.ERROR, log.Error(err), log.Data{"topic": config.ObservationProducerTopic})
			} else {
				log.Event(ctx, "kafka observation producer stopped", log.INFO, log.Data{"topic": config.ObservationProducerTopic})
			}
		}

		// cancel the timer in the shutdown context.
		cancel()

		log.Event(ctx, "graceful shutdown was successful", log.INFO)
		os.Exit(0)
	}

	// Log non-fatal errors in separate go routines
	kafkaConsumer.Channels().LogErrors(ctx, "kafka consumer error")
	kafkaObservationProducer.Channels().LogErrors(ctx, "kafka observation producer error")
	kafkaErrorProducer.Channels().LogErrors(ctx, "kafka error producer error")
	go func() {
		for {
			select {
			case err := <-errorChannel:
				log.Event(ctx, "error channel", log.ERROR, log.Error(err))
			}
		}
	}()

	// When a signal is received, shutdown gracefully
	<-signals
	log.Event(ctx, "os signal received", log.ERROR, log.Error(errors.New("os signal received")))
	shutdownGracefully()
}

// StartHealthCheck sets up the Handler, starts the healthcheck and the http server that serves health endpoint
func startHealthCheck(ctx context.Context, hc *healthcheck.HealthCheck, bindAddr string) *server.Server {
	router := mux.NewRouter()
	router.Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	httpServer := server.New(bindAddr, router)
	httpServer.HandleOSSignals = false

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Event(ctx, "http server error", log.ERROR, log.Error(err))
			hc.Stop()
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
		log.Event(nil, "error adding check for kafka consumer checker", log.ERROR, log.Error(err))
	}

	if err = hc.AddCheck("Kafka Observation Producer", kafkaObservationProducer.Checker); err != nil {
		hasErrors = true
		log.Event(nil, "error adding check for kafka observation producer checker", log.ERROR, log.Error(err))
	}

	if err = hc.AddCheck("Kafka Error Producer", kafkaErrorProducer.Checker); err != nil {
		hasErrors = true
		log.Event(nil, "error adding check for kafka error producer checker", log.ERROR, log.Error(err))
	}

	if vaultClient != nil {
		if err = hc.AddCheck("Vault", vaultClient.Checker); err != nil {
			hasErrors = true
			log.Event(nil, "error adding check for vault checker", log.ERROR, log.Error(err))
		}
	}

	for bucketName, s3 := range s3Clients {
		if err := hc.AddCheck(fmt.Sprintf("S3 bucket %s", bucketName), s3.Checker); err != nil {
			hasErrors = true
			log.Event(ctx, "error adding check for s3 client", log.ERROR, log.Error(err))
		}
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}

func checkForError(ctx context.Context, err error) {
	if err != nil {
		log.Event(ctx, "error", log.ERROR, log.Error(err))
		os.Exit(1)
	}
}
