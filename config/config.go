package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config values for the application.
type Config struct {
	AWSRegion                string        `envconfig:"AWS_REGION"`
	BindAddr                 string        `envconfig:"BIND_ADDR"`
	EncryptionDisabled       bool          `envconfig:"ENCRYPTION_DISABLED"`
	ErrorProducerTopic       string        `envconfig:"ERROR_PRODUCER_TOPIC"`
	FileConsumerGroup        string        `envconfig:"FILE_CONSUMER_GROUP"`
	FileConsumerTopic        string        `envconfig:"FILE_CONSUMER_TOPIC"`
	GracefulShutdownTimeout  time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	KafkaAddr                []string      `envconfig:"KAFKA_ADDR"                    json:"-"`
	ObservationProducerTopic string        `envconfig:"OBSERVATION_PRODUCER_TOPIC"`
	VaultAddr                string        `envconfig:"VAULT_ADDR"`
	VaultToken               string        `envconfig:"VAULT_TOKEN"                   json:"-"`
	VaultPath                string        `envconfig:"VAULT_PATH"`
	BucketNames              []string      `envconfig:"BUCKET_NAMES"                  json:"-"`
	HealthCheckInterval      time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCriticalTimeout    time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {

	cfg := &Config{
		AWSRegion:                "eu-west-1",
		BindAddr:                 ":21600",
		EncryptionDisabled:       false,
		ErrorProducerTopic:       "report-events",
		FileConsumerGroup:        "dimensions-inserted",
		FileConsumerTopic:        "dimensions-inserted",
		GracefulShutdownTimeout:  time.Second * 5,
		KafkaAddr:                []string{"localhost:9092"},
		ObservationProducerTopic: "observation-extracted",
		VaultAddr:                "http://localhost:8200",
		VaultToken:               "",
		VaultPath:                "secret/shared/psk",
		BucketNames:              []string{"dp-frontend-florence-file-uploads"},
		HealthCheckInterval:      30 * time.Second,
		HealthCriticalTimeout:    90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Config) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
