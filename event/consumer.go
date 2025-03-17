package event

import (
	"context"
	"errors"

	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/dp-reporter-client/reporter"
	"github.com/ONSdigital/log.go/v2/log"
)

// Handler represents a handler for processing a single event.
type Handler interface {
	Handle(ctx context.Context, event *DimensionsInserted) error
}

// Consumer consumes event messages.
type Consumer struct {
	Closing chan bool
	Closed  chan bool
}

// NewConsumer returns a new consumer instance.
func NewConsumer() *Consumer {
	return &Consumer{
		Closing: make(chan bool),
		Closed:  make(chan bool),
	}
}

// Consume convert them to event instances, and pass the event to the provided handler.
func (consumer *Consumer) Consume(ctx context.Context, messageConsumer kafka.IConsumerGroup, handler Handler, errorReporter reporter.ErrorReporter) {
	go func() {
		defer close(consumer.Closed)

		for {
			select {
			case message := <-messageConsumer.Channels().Upstream:

				// In the future, the context will be obtained from the kafka message
				msgCtx := context.Background()

				// Unmarshal message
				event, err := Unmarshal(message)
				if err != nil {
					log.Error(msgCtx, "message unmarshal error", err)
					message.CommitAndRelease()
					continue
				}

				logData := log.Data{"event": event}
				log.Info(msgCtx, "event received", logData)

				// Handle the message
				if err = handler.Handle(ctx, event); err != nil {
					log.Error(msgCtx, "failed to handle event", err, logData)
					if err = errorReporter.Notify(event.InstanceID, "failed to handle event", err); err != nil {
						log.Error(msgCtx, "errorReporter.Notify returned an unexpected error", err, logData)
					}
					message.CommitAndRelease()
					continue
				}

				// On success, commit and release the message
				log.Info(msgCtx, "event processed - committing message", logData)
				message.CommitAndRelease()
				log.Info(msgCtx, "message committed and kafka consumer released", logData)

			case <-consumer.Closing:
				log.Info(ctx, "closing event consumer loop")
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

	close(consumer.Closing)

	select {
	case <-consumer.Closed:
		log.Info(ctx, "successfully closed event consumer")
		return nil
	case <-ctx.Done():
		log.Info(ctx, "shutdown context time exceeded, skipping graceful shutdown of event consumer")
		return errors.New("shutdown context timed out")
	}
}

// Unmarshal converts a event instance to []byte.
func Unmarshal(message kafka.Message) (*DimensionsInserted, error) {
	var event DimensionsInserted
	err := schema.DimensionsInsertedEvent.Unmarshal(message.GetData(), &event)
	return &event, err
}
