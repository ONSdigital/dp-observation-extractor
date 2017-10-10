package event

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/dp-reporter-client/reporter"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
)

// MessageConsumer provides a generic interface for consuming []byte messages
type MessageConsumer interface {
	Incoming() chan kafka.Message
}

// Handler represents a handler for processing a single event.
type Handler interface {
	Handle(event *DimensionsInserted) error
}

// Consumer consumes event messages.
type Consumer struct {
	closing chan bool
	closed  chan bool
}

// NewConsumer returns a new consumer instance.
func NewConsumer() *Consumer {
	return &Consumer{
		closing: make(chan bool),
		closed:  make(chan bool),
	}
}

// Consume convert them to event instances, and pass the event to the provided handler.
func (consumer *Consumer) Consume(messageConsumer MessageConsumer, handler Handler, errorReporter reporter.ErrorReporter) {

	go func() {
		defer close(consumer.closed)

		for {
			select {
			case message := <-messageConsumer.Incoming():

				event, err := Unmarshal(message)
				if err != nil {
					log.Error(err, log.Data{"message": "failed to unmarshal event"})
					continue
				}

				logData := log.Data{"event": event}
				log.Debug("event received", logData)

				err = handler.Handle(event)
				if err != nil {
					log.ErrorC("failed to handle event", err, logData)
					if err := errorReporter.Notify(event.InstanceID, "failed to handle event", err); err != nil {
						log.ErrorC("errorReporter.Notify returned an unexpected error", err, logData)
					}
					continue
				}

				log.Debug("event processed - committing message", logData)
				message.Commit()
				log.Debug("message committed", logData)

			case <-consumer.closing:
				log.Info("closing event consumer loop", nil)
				return
			}
		}

	}()
}

// Close safely closes the consumer and releases all resources
func (consumer *Consumer) Close(ctx context.Context) (err error) {

	if ctx == nil {
		ctx = context.Background()
	}

	close(consumer.closing)

	select {
	case <-consumer.closed:
		log.Info("successfully closed event consumer", nil)
		return nil
	case <-ctx.Done():
		log.Info("shutdown context time exceeded, skipping graceful shutdown of event consumer", nil)
		return errors.New("shutdown context timed out")
	}
}

// Unmarshal converts a event instance to []byte.
func Unmarshal(message kafka.Message) (*DimensionsInserted, error) {
	var event DimensionsInserted
	err := schema.DimensionsInsertedEvent.Unmarshal(message.GetData(), &event)
	return &event, err
}
