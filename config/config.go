package config

import "github.com/kelseyhightower/envconfig"

// Config values for the application.
type Config struct {
	BindAddr                 string `envconfig:"BIND_ADDR"`
	KafkaAddr                string `envconfig:"KAFKA_ADDR"`
	FileConsumerGroup        string `envconfig:"FILE_CONSUMER_GROUP"`
	FileConsumerTopic        string `envconfig:"FILE_CONSUMER_TOPIC"`
	AWSRegion                string `envconfig:"AWS_REGION"`
	ObservationProducerTopic string `envconfig:"OBSERVATION_PRODUCER_TOPIC"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		BindAddr:                 ":21600",
		KafkaAddr:                "localhost:9092",
		FileConsumerGroup:        "dimensions-inserted",
		FileConsumerTopic:        "dimensions-inserted",
		AWSRegion:                "eu-west-1",
		ObservationProducerTopic: "observation-extracted",
	}

	return cfg, envconfig.Process("", cfg)
}
