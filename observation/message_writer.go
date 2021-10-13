package observation

import (
	"context"

	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/log.go/v2/log"
)

// MessageWriter writes observations as messages
type MessageWriter struct {
	messageProducer MessageProducer
}

// MessageProducer dependency that writes messages
type MessageProducer interface {
	Channels() *kafka.ProducerChannels
}

// NewMessageWriter returns a new observation message writer.
func NewMessageWriter(messageProducer MessageProducer) *MessageWriter {
	return &MessageWriter{
		messageProducer: messageProducer,
	}
}

// WriteAll observations as messages from the given observation reader.
func (messageWriter MessageWriter) WriteAll(ctx context.Context, reader Reader, instanceID string) {

	observation, readErr := reader.Read()

	for readErr == nil {

		extractedEvent := ExtractedEvent{
			InstanceID: instanceID,
			Row:        observation.Row,
			RowIndex:   observation.RowIndex,
		}

		bytes, err := schema.ObservationExtractedEvent.Marshal(extractedEvent)
		if err != nil {
			log.Error(ctx, "", err, log.Data{
				"schema": "failed to marshal observation extracted event",
				"event":  extractedEvent})
		}

		messageWriter.messageProducer.Channels().Output <- bytes

		observation, readErr = reader.Read()
	}

	log.Info(ctx, "all observations extracted", log.Data{"instanceID": instanceID})
}

// Marshal converts the given observationExtractedEvent to a []byte.
func Marshal(extractedEvent ExtractedEvent) ([]byte, error) {
	bytes, err := schema.ObservationExtractedEvent.Marshal(extractedEvent)
	return bytes, err
}
