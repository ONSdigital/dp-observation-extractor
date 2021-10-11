package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// KafkaTLSProtocolFlag informs service to use TLS protocol for kafka
const KafkaTLSProtocolFlag = "TLS"

// Config values for the application.
type Config struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	AWSRegion               string        `envconfig:"AWS_REGION"`
	BucketNames             []string      `envconfig:"BUCKET_NAMES"                   json:"-"`
	EncryptionDisabled      bool          `envconfig:"ENCRYPTION_DISABLED"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval     time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCriticalTimeout   time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	KafkaConfig             KafkaConfig
	VaultAddr               string `envconfig:"VAULT_ADDR"`
	VaultToken              string `envconfig:"VAULT_TOKEN"                           json:"-"`
	VaultPath               string `envconfig:"VAULT_PATH"`
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Brokers                  []string `envconfig:"KAFKA_ADDR"                         json:"-"`
	Version                  string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest             bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	SecProtocol              string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts               string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientKey             string   `envconfig:"KAFKA_SEC_CLIENT_KEY"               json:"-"`
	SecClientCert            string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecSkipVerify            bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	ErrorProducerTopic       string   `envconfig:"ERROR_PRODUCER_TOPIC"`
	FileConsumerGroup        string   `envconfig:"FILE_CONSUMER_GROUP"`
	FileConsumerTopic        string   `envconfig:"FILE_CONSUMER_TOPIC"`
	ObservationProducerTopic string   `envconfig:"OBSERVATION_PRODUCER_TOPIC"`
}

func getDefaultConfig() *Config {
	return &Config{
		BindAddr:                ":21600",
		AWSRegion:               "eu-west-1",
		BucketNames:             []string{"dp-frontend-florence-file-uploads"},
		EncryptionDisabled:      false,
		GracefulShutdownTimeout: time.Second * 5,
		HealthCheckInterval:     30 * time.Second,
		HealthCriticalTimeout:   90 * time.Second,
		KafkaConfig: KafkaConfig{
			Brokers:                  []string{"localhost:9092"},
			Version:                  "1.0.2",
			OffsetOldest:             true,
			SecProtocol:              "",
			SecCACerts:               "",
			SecClientCert:            "",
			SecClientKey:             "",
			SecSkipVerify:            false,
			ErrorProducerTopic:       "report-events",
			FileConsumerGroup:        "dimensions-inserted",
			FileConsumerTopic:        "dimensions-inserted",
			ObservationProducerTopic: "observation-extracted",
		},
		VaultAddr:  "http://localhost:8200",
		VaultToken: "",
		VaultPath:  "secret/shared/psk",
	}
}

// Get the configuration values from the environment or provide the defaults.
func Get() (*Config, error) {
	cfg := getDefaultConfig()

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	if errs := cfg.KafkaConfig.validate(); len(errs) != 0 {
		return nil, fmt.Errorf("kafka config validation errors: %v", strings.Join(errs, ", "))
	}

	return cfg, nil
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Config) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
