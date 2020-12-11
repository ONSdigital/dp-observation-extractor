package event_test

import (
	"context"

	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/event/eventtest"
	"github.com/ONSdigital/dp-observation-extractor/schema"

	"errors"
	"testing"

	"github.com/ONSdigital/dp-reporter-client/reporter/reportertest"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

func TestConsume_UnmarshallError(t *testing.T) {
	Convey("Given an event consumer with an invalid schema and a valid schema", t, func(c C) {

		reporter := reportertest.NewImportErrorReporterMock(nil)
		messageConsumer := kafkatest.NewMessageConsumer(true)
		handler := eventtest.NewEventHandler(nil)

		expectedEvent := getExampleEvent()

		messageIncorrect := kafkatest.NewMessage([]byte("invalid schema"), 0)
		messageCorrect := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)

		// Make sure the invalid schema is sent first
		go func() {
			messageConsumer.Channels().Upstream <- messageIncorrect
			messageConsumer.Channels().Upstream <- messageCorrect
		}()

		Convey("When consume messages is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			// Wait for handler to receive message, and message to be successfully released
			<-handler.ChHandle
			<-messageCorrect.UpstreamDone()
			<-messageIncorrect.UpstreamDone()

			Convey("Only the valid event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})

			Convey("And errorHandler is never called", func() {
				So(len(reporter.NotifyCalls()), ShouldEqual, 0)
			})
		})
	})
}

func TestConsumer_HandlerError(t *testing.T) {
	Convey("Given an event consumer with a valid schema", t, func(c C) {

		reporter := reportertest.NewImportErrorReporterMock(nil)
		messageConsumer := kafkatest.NewMessageConsumer(true)

		handlerErr := errors.New("handler error")
		handler := eventtest.NewEventHandler(handlerErr)

		expectedEvent := getExampleEvent()

		messageConsumer.Channels().Upstream <- kafkatest.NewMessage(marshal(*expectedEvent, c), 0)

		Convey("When consume is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			waitEventsAndCloseHandler(ctx, consumer, handler, 1)

			Convey("An event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)

				Convey("Then the returned handler error is passed to the error handler", func() {
					So(len(reporter.NotifyCalls()), ShouldEqual, 1)
					So(reporter.NotifyCalls()[0].ID, ShouldEqual, expectedEvent.InstanceID)
					So(reporter.NotifyCalls()[0].ErrContext, ShouldEqual, "failed to handle event")
					So(reporter.NotifyCalls()[0].Err, ShouldResemble, handlerErr)
				})
			})
		})
	})
}

func TestConsume(t *testing.T) {

	Convey("Given an event consumer with a valid schema", t, func(c C) {

		reporter := reportertest.NewImportErrorReporterMock(nil)
		messageConsumer := kafkatest.NewMessageConsumer(true)
		handler := eventtest.NewEventHandler(nil)

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)
		messageConsumer.Channels().Upstream <- message

		Convey("When consume is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			waitEventsAndCloseHandler(ctx, consumer, handler, 1)

			Convey("A event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})

			Convey("The message is committed, and consumer group is released", func() {
				So(len(message.CommitAndReleaseCalls()), ShouldEqual, 1)
			})

			Convey("And errorHandler is never called", func() {
				So(len(reporter.NotifyCalls()), ShouldEqual, 0)
			})
		})
	})
}

func TestToEvent(t *testing.T) {

	Convey("Given a event schema encoded using avro", t, func(c C) {

		expectedEvent := getExampleEvent()
		message := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)

		Convey("When the expectedEvent is unmarshalled", func() {

			event, err := event.Unmarshal(message)

			Convey("The expectedEvent has the expected values", func() {
				So(err, ShouldBeNil)
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})
		})
	})
}

// marshal helper method to marshal a event into a []byte
func marshal(event event.DimensionsInserted, c C) []byte {
	bytes, err := schema.DimensionsInsertedEvent.Marshal(event)
	c.So(err, ShouldBeNil)
	return bytes
}

func getExampleEvent() *event.DimensionsInserted {
	expectedEvent := &event.DimensionsInserted{
		InstanceID: "1234",
		FileURL:    "s3://some-bucket/some-file",
	}
	return expectedEvent
}

// waitEventsAndCloseHandler waits for a number or events to be sent to the handler, and then closes it,
// blocking until closed channel is closed, in order to prevent race conditions
func waitEventsAndCloseHandler(ctx context.Context, consumer *event.Consumer, handler *eventtest.EventHandler, nMsg int) {

	// Wait for handler to receive nMsg messages
	for i := 0; i < nMsg; i++ {
		<-handler.ChHandle
	}

	// Close the consumer, and wait for it to be closed
	err := consumer.Close(ctx)
	So(err, ShouldBeNil)
	<-consumer.Closed
}
