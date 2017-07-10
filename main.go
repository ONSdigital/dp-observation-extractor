package main

import (
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/kafka"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/s3"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	log.Debug("Loaded config", log.Data{"config": cfg})

	s3, err := s3.New(cfg.AWSRegion)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	kafkaConsumer := kafka.Consumer{}
	kafkaProducer := kafka.Producer{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		// gracefully dispose resources

		log.Debug("Graceful shutdown of  was successful.", nil)
		os.Exit(0)
	}()

	observationWriter := observation.NewMessageWriter(kafkaProducer)
	requestHandler := request.NewCSVHandler(s3, observationWriter)

	request.Consume(kafkaConsumer, requestHandler)
}
