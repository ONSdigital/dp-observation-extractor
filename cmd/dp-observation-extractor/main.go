package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

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

const chunkSize = 5 * 1024 * 1024

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

	// TODO add new healthcheck
	router := mux.NewRouter()
	// router.Path("/healthcheck").HandlerFunc(healthcheck.Handler)
	httpServer := server.New(config.BindAddr, router)

	// Disable auto handling of os signals by the HTTP server. This is handled
	// in the service so we can gracefully shutdown resources other than just
	// the HTTP server.
	httpServer.HandleOSSignals = false

	go func() {
		log.Event(ctx, "starting http server", log.INFO, log.Data{"bind_addr": config.BindAddr})
		if err = httpServer.ListenAndServe(); err != nil {
			errorChannel <- err
		}
	}()

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

	// var cryptoClient event.CryptoClient
	var vaultClient event.VaultClient
	if !config.EncryptionDisabled {
		// cryptoClient = s3crypto.New(sess, &s3crypto.Config{HasUserDefinedPSK: true, MultipartChunkSize: chunkSize})
		vaultClient, err = vault.CreateClient(config.VaultToken, config.VaultAddr, 3)
		checkForError(ctx, err)
	}
	// client := s3.New(sess)

	// bucket := "myBucket"
	// s3cli, err := s3client.NewClient(config.AWSRegion, bucket, !config.EncryptionDisabled)
	// checkForError(ctx, err)

	eventHandler := event.NewCSVHandler(sess, s3Clients, vaultClient, observationWriter, config.VaultPath)

	errorReporter, err := reporter.NewImportErrorReporter(kafkaErrorProducer, log.Namespace)
	checkForError(ctx, err)

	eventConsumer := event.NewConsumer()
	eventConsumer.Consume(ctx, kafkaConsumer, eventHandler, errorReporter)

	shutdownGracefully := func() {

		ctx, cancel := context.WithTimeout(ctx, config.GracefulShutdownTimeout)

		// gracefully dispose resources
		err = eventConsumer.Close(ctx)
		if err != nil {
			log.Event(ctx, "error closing event consumer", log.ERROR, log.Error(err))
		}

		err = kafkaConsumer.Close(ctx)
		if err != nil {
			log.Event(ctx, "error closing kafka consumer", log.ERROR, log.Error(err))
		}

		err = kafkaErrorProducer.Close(ctx)
		if err != nil {
			log.Event(ctx, "error closing kafka error producer", log.ERROR, log.Error(err))
		}

		err = kafkaObservationProducer.Close(ctx)
		if err != nil {
			log.Event(ctx, "error closing kafka observation producer", log.ERROR, log.Error(err))
		}

		err = httpServer.Shutdown(ctx)
		if err != nil {
			log.Event(ctx, "error shutting down http server ", log.ERROR, log.Error(err))
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

func checkForError(ctx context.Context, err error) {
	if err != nil {
		log.Event(ctx, "error", log.ERROR, log.Error(err))
		os.Exit(1)
	}
}
