package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/go-ns/errorhandler"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/s3"
)

func main() {

	log.Namespace = "dp-observation-extractor"

	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	log.Debug("loaded config", log.Data{"config": config})

	kafkaBrokers := []string{config.KafkaAddr}

	s3, err := s3.New(config.AWSRegion)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	kafkaConsumer, err := kafka.NewConsumerGroup(kafkaBrokers,
		config.FileConsumerTopic,
		config.FileConsumerGroup,
		kafka.OffsetNewest)
	if err != nil {
		log.Error(err, log.Data{"message": "failed to create kafka consumer"})
		os.Exit(1)
	}

	kafkaProducer := kafka.NewProducer(kafkaBrokers, config.ObservationProducerTopic, 0)
	kafkaErrorProducer := kafka.NewProducer(kafkaBrokers, config.ErrorProducerTopic, 0)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	errorHandler := errorhandler.NewKafkaHandler(kafkaErrorProducer)

	go func() {
		<-signals

		// gracefully dispose resources
		kafkaConsumer.Closer() <- true
		kafkaProducer.Closer() <- true
		kafkaErrorProducer.Closer() <- true

		log.Debug("graceful shutdown was successful", nil)
		os.Exit(0)
	}()

	observationWriter := observation.NewMessageWriter(kafkaProducer, errorHandler)
	eventHandler := event.NewCSVHandler(s3, observationWriter, errorHandler)

	event.Consume(kafkaConsumer, eventHandler)
}
