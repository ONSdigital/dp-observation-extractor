package initialise

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	s3client "github.com/ONSdigital/dp-s3/v3"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ExternalServiceList represents a list of services
type ExternalServiceList struct {
	Consumer              bool
	ObservationProducer   bool
	ErrorReporterProducer bool
	Vault                 bool
	HealthCheck           bool
	S3Clients             bool
}

// KafkaProducerName represents a type for kafka producer name used by iota constants
type KafkaProducerName int

// Possible names of Kafka Producers
const (
	Observation = iota
	ErrorReporter
)

var kafkaProducerNames = []string{"Observation", "ErrorReporter"}

// Values of the kafka producers names
func (k KafkaProducerName) String() string {
	return kafkaProducerNames[k]
}

// GetConsumer returns a kafka consumer, which might not be initialised
func (e *ExternalServiceList) GetConsumer(ctx context.Context, kafkaConfig *config.KafkaConfig) (*kafka.ConsumerGroup, error) {
	kafkaOffset := kafka.OffsetNewest

	if kafkaConfig.OffsetOldest {
		kafkaOffset = kafka.OffsetOldest
	}

	cgConfig := &kafka.ConsumerGroupConfig{
		Offset:       &kafkaOffset,
		KafkaVersion: &kafkaConfig.Version,
	}
	if kafkaConfig.SecProtocol == config.KafkaTLSProtocolFlag {
		cgConfig.SecurityConfig = kafka.GetSecurityConfig(
			kafkaConfig.SecCACerts,
			kafkaConfig.SecClientCert,
			kafkaConfig.SecClientKey,
			kafkaConfig.SecSkipVerify,
		)
	}

	kafkaConsumer, err := kafka.NewConsumerGroup(
		ctx,
		kafkaConfig.Brokers,
		kafkaConfig.FileConsumerTopic,
		kafkaConfig.FileConsumerGroup,
		kafka.CreateConsumerGroupChannels(1),
		cgConfig,
	)
	if err != nil {
		return kafkaConsumer, err
	}

	e.Consumer = true

	return kafkaConsumer, nil
}

// GetProducer returns a kafka producer, which might not be initialised
func (e *ExternalServiceList) GetProducer(ctx context.Context, kafkaConfig *config.KafkaConfig, topic string, name KafkaProducerName) (kafkaProducer *kafka.Producer, err error) {
	pChannels := kafka.CreateProducerChannels()

	pConfig := &kafka.ProducerConfig{
		KafkaVersion: &kafkaConfig.Version,
	}
	if kafkaConfig.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			kafkaConfig.SecCACerts,
			kafkaConfig.SecClientCert,
			kafkaConfig.SecClientKey,
			kafkaConfig.SecSkipVerify,
		)
	}

	producer, err := kafka.NewProducer(ctx, kafkaConfig.Brokers, topic, pChannels, pConfig)
	if err != nil {
		log.Fatal(ctx, "new kafka producer returned an error", err, log.Data{"topic": topic})
		return nil, err
	}

	switch {
	case name == Observation:
		e.ObservationProducer = true
	case name == ErrorReporter:
		e.ErrorReporterProducer = true
	default:
		return nil, fmt.Errorf("kafka producer name not recognised: '%s'. valid names: %v", name.String(), kafkaProducerNames)
	}

	return producer, nil
}

// GetS3Clients returns a map of AWS S3 clients corresponding to the list of BucketNames
// and the AWS region provided in the configuration. If encryption is enabled, the s3clients will be cryptoclients.
func (e *ExternalServiceList) GetS3Clients(ctx context.Context, cfg *config.Config) (*aws.Config, map[string]event.S3Client, error) {
	if cfg.LocalstackHost != "" {
		awsConfig, err := awsConfig.LoadDefaultConfig(ctx,
			awsConfig.WithRegion(cfg.AWSRegion),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		)
		if err != nil {
			return nil, nil, err
		}

		s3Clients := make(map[string]event.S3Client)
		for _, bucketName := range cfg.BucketNames {
			s3Clients[bucketName] = s3client.NewClientWithConfig(bucketName, awsConfig, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(cfg.LocalstackHost)
				o.UsePathStyle = true
			})
		}
		e.S3Clients = true

		return &awsConfig, s3Clients, nil
	}

	// establish AWS config
	awsConfig, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		return nil, nil, err
	}

	// create S3 clients for expected bucket names, so that they can be health-checked
	s3Clients := make(map[string]event.S3Client)
	for _, bucketName := range cfg.BucketNames {
		s3Clients[bucketName] = s3client.NewClientWithConfig(bucketName, awsConfig)
	}
	e.S3Clients = true

	return &awsConfig, s3Clients, nil
}

// GetHealthChecker creates a new healthcheck object
func (e *ExternalServiceList) GetHealthChecker(ctx context.Context, buildTime, gitCommit, version string, cfg *config.Config) (*healthcheck.HealthCheck, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "failed to create versionInfo for healthcheck", err)
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCriticalTimeout, cfg.HealthCheckInterval)
	e.HealthCheck = true

	return &hc, nil
}
