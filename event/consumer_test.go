package event_test

import (
	"context"

	"github.com/ONSdigital/dp-kafka/kafkatest"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/event/eventtest"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/log.go/log"

	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-reporter-client/reporter/reportertest"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

func TestConsume_UnmarshallError(t *testing.T) {
	Convey("Given an event consumer with an invalid schema and a valid schema", t, func(c C) {

		reporter := reportertest.NewImportErrorReporterMock(nil)
		// messages := make(chan kafka.Message, 2)
		messageConsumer := kafkatest.NewMessageConsumer(true)
		// messageConsumer := kafkatest.NewMessageConsumer(messages)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		go func() {
			messageConsumer.Channels().Upstream <- kafkatest.NewMessage([]byte("invalid schema"), 0)
			messageConsumer.Channels().Upstream <- kafkatest.NewMessage(marshal(*expectedEvent, c), 0)
		}()

		Convey("When consume messages is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			waitForEventsToBeSentToHandler(ctx, handler)

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
		// messages := make(chan kafka.Message, 1)

		messageConsumer := kafkatest.NewMessageConsumer(true)

		handlerErr := errors.New("handler error")
		handler := &eventtest.EventHandler{
			Events: nil,
			Error:  handlerErr,
		}

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)
		messageConsumer.Channels().Upstream <- message

		Convey("When consume is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			waitForEventsToBeSentToHandler(ctx, handler)

			Convey("A event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})

			Convey("Then the returned handler error is passed to the error handler", func() {
				So(len(reporter.NotifyCalls()), ShouldEqual, 1)
				So(reporter.NotifyCalls()[0].ID, ShouldEqual, expectedEvent.InstanceID)
				So(reporter.NotifyCalls()[0].ErrContext, ShouldEqual, "failed to handle event")
				So(reporter.NotifyCalls()[0].Err, ShouldResemble, handlerErr)
			})
		})
	})
}

func TestConsume(t *testing.T) {

	Convey("Given an event consumer with a valid schema", t, func(c C) {

		reporter := reportertest.NewImportErrorReporterMock(nil)
		// messages := make(chan kafka.Message, 1)
		messageConsumer := kafkatest.NewMessageConsumer(true)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)
		messageConsumer.Channels().Upstream <- message

		Convey("When consume is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(ctx, messageConsumer, handler, reporter)

			waitForEventsToBeSentToHandler(ctx, handler)

			Convey("A event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})

			Convey("The message is committed, and consumer group is released", func() {
				So(len(messageConsumer.CommitAndReleaseCalls()), ShouldEqual, 1)
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

func TestClose(t *testing.T) {

	Convey("Given a consumer", t, func(c C) {

		// messages := make(chan kafka.Message, 1)
		messageConsumer := kafkatest.NewMessageConsumer(true)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent, c), 0)
		messageConsumer.Channels().Upstream <- message

		consumer := event.NewConsumer()
		consumer.Consume(ctx, messageConsumer, handler, nil)

		Convey("When close is called", func() {

			err := consumer.Close(context.Background())

			Convey("The expected event is sent to the handler", func() {
				So(err, ShouldBeNil)
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

func waitForEventsToBeSentToHandler(ctx context.Context, eventHandler *eventtest.EventHandler) {

	start := time.Now()
	timeout := start.Add(time.Millisecond * 500)
	for {
		if len(eventHandler.Events) > 0 {
			log.Event(ctx, "events have been sent to the handler", log.INFO)
			break
		}

		if time.Now().After(timeout) {
			log.Event(ctx, "timeout hit", log.INFO)
			break
		}

		time.Sleep(time.Millisecond * 10)
	}
}
