package initialise

import (
	"context"
	"fmt"

	kafka "github.com/ONSdigital/dp-kafka"
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/log.go/log"
)

// ExternalServiceList represents a list of services
type ExternalServiceList struct {
	Consumer              bool
	ObservationProducer   bool
	ErrorReporterProducer bool
	Vault                 bool
	HealthCheck           bool
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
func (e *ExternalServiceList) GetConsumer(ctx context.Context, cfg *config.Config) (kafkaConsumer *kafka.ConsumerGroup, err error) {
	kafkaConsumer, err = kafka.NewConsumerGroup(
		ctx,
		cfg.KafkaAddr,
		cfg.FileConsumerTopic,
		cfg.FileConsumerGroup,
		kafka.OffsetNewest,
		true,
		kafka.CreateConsumerGroupChannels(true),
	)
	if err != nil {
		return
	}

	e.Consumer = true

	return
}

// GetProducer returns a kafka producer, which might not be initialised
func (e *ExternalServiceList) GetProducer(ctx context.Context, brokers []string, topic string, name KafkaProducerName) (kafkaProducer *kafka.Producer, err error) {
	producer, err := kafka.NewProducer(ctx, brokers, topic, 0, kafka.CreateProducerChannels())
	if err != nil {
		log.Event(ctx, "new kafka producer returned an error", log.FATAL, log.Error(err), log.Data{"topic": topic})
		return nil, err
	}

	switch {
	case name == Observation:
		e.ObservationProducer = true
	case name == ErrorReporter:
		e.ErrorReporterProducer = true
	default:
		err = fmt.Errorf("kafka producer name not recognised: '%s'. valid names: %v", name.String(), kafkaProducerNames)
	}

	return producer, nil
}
