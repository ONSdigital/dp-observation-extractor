package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

// Config values for the application.
type Config struct {
	BindAddr                 string        `envconfig:"BIND_ADDR"`
	KafkaAddr                []string      `envconfig:"KAFKA_ADDR"`
	FileConsumerGroup        string        `envconfig:"FILE_CONSUMER_GROUP"`
	FileConsumerTopic        string        `envconfig:"FILE_CONSUMER_TOPIC"`
	AWSRegion                string        `envconfig:"AWS_REGION"`
	ObservationProducerTopic string        `envconfig:"OBSERVATION_PRODUCER_TOPIC"`
	ErrorProducerTopic       string        `envconfig:"ERROR_PRODUCER_TOPIC"`
	GracefulShutdownTimeout  time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		BindAddr:                 ":21600",
		KafkaAddr:                []string{"localhost:9092"},
		FileConsumerGroup:        "dimensions-inserted",
		FileConsumerTopic:        "dimensions-inserted",
		ErrorProducerTopic:       "import-error",
		AWSRegion:                "eu-west-1",
		ObservationProducerTopic: "observation-extracted",
		GracefulShutdownTimeout:  time.Second * 5,
	}

	return cfg, envconfig.Process("", cfg)
}
