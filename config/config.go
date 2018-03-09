package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config values for the application.
type Config struct {
	AWSPrivateKey            string        `envconfig:"RSA_PRIVATE_KEY"               json:"-"`
	AWSRegion                string        `envconfig:"AWS_REGION"`
	BindAddr                 string        `envconfig:"BIND_ADDR"`
	EncryptionDisabled       bool          `envconfig:"ENCRYPTION_DISABLED"`
	ErrorProducerTopic       string        `envconfig:"ERROR_PRODUCER_TOPIC"`
	FileConsumerGroup        string        `envconfig:"FILE_CONSUMER_GROUP"`
	FileConsumerTopic        string        `envconfig:"FILE_CONSUMER_TOPIC"`
	GracefulShutdownTimeout  time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	KafkaAddr                []string      `envconfig:"KAFKA_ADDR"                    json:"-"`
	ObservationProducerTopic string        `envconfig:"OBSERVATION_PRODUCER_TOPIC"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		AWSRegion:                "eu-west-1",
		BindAddr:                 ":21600",
		EncryptionDisabled:       true,
		ErrorProducerTopic:       "report-events",
		FileConsumerGroup:        "dimensions-inserted",
		FileConsumerTopic:        "dimensions-inserted",
		GracefulShutdownTimeout:  time.Second * 5,
		KafkaAddr:                []string{"localhost:9092"},
		ObservationProducerTopic: "observation-extracted",
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Config) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
