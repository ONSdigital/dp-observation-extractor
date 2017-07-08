package main

import (
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"os/signal"
	"syscall"
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/go-ns/s3"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/kafka"
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

		//if err := extractor.Close(); err != nil {
		//	log.Debug("Failed to shutdown gracefully.", nil)
		//	log.Error(err, nil)
		//	os.Exit(1)
		//}
		log.Debug("Graceful shutdown of  was successful.", nil)
		os.Exit(0)
	}()


	observationWriter := observation.NewMessageWriter(kafkaProducer)
	requestHandler := request.NewCSVHandler(s3, observationWriter)

	request.Consume(kafkaConsumer, requestHandler)
}
