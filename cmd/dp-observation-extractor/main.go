package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-reporter-client/reporter"
	"github.com/ONSdigital/go-ns/handlers/healthcheck"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/go-ns/vault"
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

func main() {

	log.Namespace = "dp-observation-extractor"

	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	// sensitive fields are omitted from config.String().
	log.Info("config on startup", log.Data{"config": config})

	// a channel used to signal a graceful exit is required.
	errorChannel := make(chan error)

	router := mux.NewRouter()
	router.Path("/healthcheck").HandlerFunc(healthcheck.Handler)
	httpServer := server.New(config.BindAddr, router)

	// Disable auto handling of os signals by the HTTP server. This is handled
	// in the service so we can gracefully shutdown resources other than just
	// the HTTP server.
	httpServer.HandleOSSignals = false

	go func() {
		log.Debug("starting http server", log.Data{"bind_addr": config.BindAddr})
		if err = httpServer.ListenAndServe(); err != nil {
			errorChannel <- err
		}
	}()

	sess, err := session.NewSession(&aws.Config{Region: &config.AWSRegion})
	checkForError(err)

	kafkaConsumer, err := kafka.NewConsumerGroup(config.KafkaAddr,
		config.FileConsumerTopic,
		config.FileConsumerGroup,
		kafka.OffsetNewest)
	checkForError(err)

	kafkaObservationProducer, err := kafka.NewProducer(config.KafkaAddr, config.ObservationProducerTopic, 0)
	checkForError(err)

	kafkaErrorProducer, err := kafka.NewProducer(config.KafkaAddr, config.ErrorProducerTopic, 0)
	checkForError(err)

	observationWriter := observation.NewMessageWriter(kafkaObservationProducer)

	var cryptoClient event.CryptoClient
	var vaultClient event.VaultClient
	if !config.EncryptionDisabled {
		cryptoClient = s3crypto.New(sess, &s3crypto.Config{HasUserDefinedPSK: true})
		vaultClient, err = vault.CreateVaultClient(config.VaultToken, config.VaultAddr, 3)
		checkForError(err)
	}
	client := s3.New(sess)
	eventHandler := event.NewCSVHandler(client, cryptoClient, vaultClient, observationWriter, config.VaultPath)

	errorReporter, err := reporter.NewImportErrorReporter(kafkaErrorProducer, log.Namespace)
	checkForError(err)

	eventConsumer := event.NewConsumer()
	eventConsumer.Consume(kafkaConsumer, eventHandler, errorReporter)

	shutdownGracefully := func() {

		ctx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)

		// gracefully dispose resources
		err = eventConsumer.Close(ctx)
		if err != nil {
			log.Error(err, nil)
		}

		err = kafkaConsumer.Close(ctx)
		if err != nil {
			log.Error(err, nil)
		}

		err = kafkaErrorProducer.Close(ctx)
		if err != nil {
			log.Error(err, nil)
		}

		err = kafkaObservationProducer.Close(ctx)
		if err != nil {
			log.Error(err, nil)
		}

		err = httpServer.Shutdown(ctx)
		if err != nil {
			log.Error(err, nil)
		}

		// cancel the timer in the shutdown context.
		cancel()

		log.Debug("graceful shutdown was successful", nil)
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case err := <-kafkaConsumer.Errors():
			log.ErrorC("kafka consumer", err, nil)
			shutdownGracefully()
		case err := <-kafkaObservationProducer.Errors():
			log.ErrorC("kafka result producer", err, nil)
			shutdownGracefully()
		case err := <-kafkaErrorProducer.Errors():
			log.ErrorC("kafka error producer", err, nil)
			shutdownGracefully()
		case err := <-errorChannel:
			log.ErrorC("error channel", err, nil)
			shutdownGracefully()
		case <-signals:
			log.Error(errors.New("os signal received"), nil)
			shutdownGracefully()
		}
	}
}

func checkForError(err error) {
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
}
