package config

import (
	"github.com/ian-kent/gofigure"
)

// Config values for the application.
type Config struct {
	BindAddr                 string `env:"BIND_ADDR" flag:"bind-addr" flagDesc:"The port to bind to"`
	KafkaAddr                string `env:"KAFKA_ADDR" flag:"kafka-addr" flagDesc:"The address of the Kafka instance"`
	FileConsumerGroup        string `env:"FILE_CONSUMER_GROUP" flag:"file-consumer-group" flagDesc:"The Kafka consumer group to consume file messages from"`
	FileConsumerTopic        string `env:"FILE_CONSUMER_TOPIC" flag:"file-consumer-topic" flagDesc:"The Kafka topic to consume file messages from"`
	AWSRegion                string `env:"AWS_REGION" flag:"aws-region" flagDesc:"The AWS region to use"`
	ObservationProducerTopic string `env:"OBSERVATION_PRODUCER_TOPIC" flag:"observation-producer-topic" flagDesc:"The Kafka topic to send the observation messages to"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := Config{
		BindAddr:                 ":21600",
		KafkaAddr:                "http://localhost:9092",
		FileConsumerGroup:        "dimensions-inserted",
		FileConsumerTopic:        "dimensions-inserted",
		AWSRegion:                "eu-west-1",
		ObservationProducerTopic: "observation-io-extracted",
	}

	err := gofigure.Gofigure(&cfg)

	return &cfg, err
}
