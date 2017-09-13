package observation

import (
	"github.com/ONSdigital/dp-observation-extractor/errors"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/go-ns/log"
)

// MessageWriter writes observations as messages
type MessageWriter struct {
	messageProducer MessageProducer
	errorHandler    errors.Handler
}

// MessageProducer dependency that writes messages
type MessageProducer interface {
	Output() chan []byte
	Closer() chan bool
}

// NewMessageWriter returns a new observation message writer.
func NewMessageWriter(messageProducer MessageProducer, errorHandler errors.Handler) *MessageWriter {
	return &MessageWriter{
		messageProducer: messageProducer,
		errorHandler:    errorHandler,
	}
}

// WriteAll observations as messages from the given observation reader.
func (messageWriter MessageWriter) WriteAll(reader Reader, instanceID string) {

	observation, readErr := reader.Read()

	for readErr == nil {

		extractedEvent := ExtractedEvent{
			InstanceID: instanceID,
			Row:        observation.Row,
		}

		bytes, err := schema.ObservationExtractedEvent.Marshal(extractedEvent)
		if err != nil {
			messageWriter.errorHandler.Handle(instanceID, err)
			log.Error(err, log.Data{
				"schema": "failed to marshal observation extracted event",
				"event":  extractedEvent})
		}

		messageWriter.messageProducer.Output() <- bytes

		observation, readErr = reader.Read()
	}

	log.Debug("all observations extracted", log.Data{"instanceID": instanceID})
}

// Marshal converts the given observationExtractedEvent to a []byte.
func Marshal(extractedEvent ExtractedEvent) ([]byte, error) {
	bytes, err := schema.ObservationExtractedEvent.Marshal(extractedEvent)
	return bytes, err
}
