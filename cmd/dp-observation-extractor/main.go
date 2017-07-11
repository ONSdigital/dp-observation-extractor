package main

import (
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/kafka"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/s3"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	log.Debug("loaded config", log.Data{"config": config})

	s3, err := s3.New(config.AWSRegion)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	kafkaBrokers := []string{config.KafkaAddr}

	kafkaConsumer, err := kafka.NewConsumerGroup(config.FileConsumerTopic, config.FileConsumerGroup)
	if err != nil {
		log.Error(err, log.Data{"message": "failed to create kafka consumer"})
		os.Exit(1)
	}

	kafkaProducer := kafka.NewProducer(kafkaBrokers, config.ObservationProducerTopic, 0)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		// gracefully dispose resources
		kafkaConsumer.Closer() <- true
		kafkaProducer.Closer() <- true

		log.Debug("graceful shutdown was successful", nil)
		os.Exit(0)
	}()

	observationWriter := observation.NewMessageWriter(kafkaProducer)
	eventHandler := event.NewCSVHandler(s3, observationWriter)

	event.Consume(kafkaConsumer, eventHandler)
}
