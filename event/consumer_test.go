package event_test

import (
	"context"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/event/eventtest"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/kafka/kafkatest"
	"github.com/ONSdigital/go-ns/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestConsume_UnmarshallError(t *testing.T) {
	Convey("Given an event consumer with an invalid schema and a valid schema", t, func() {

		messages := make(chan kafka.Message, 2)
		messageConsumer := kafkatest.NewMessageConsumer(messages)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		messages <- kafkatest.NewMessage([]byte("invalid schema"))
		messages <- kafkatest.NewMessage(marshal(*expectedEvent))

		Convey("When consume messages is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(messageConsumer, handler)

			waitForEventsToBeSentToHandler(handler)

			Convey("Only the valid event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})
		})
	})
}

func TestConsume(t *testing.T) {

	Convey("Given an event consumer with a valid schema", t, func() {

		messages := make(chan kafka.Message, 1)
		messageConsumer := kafkatest.NewMessageConsumer(messages)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent))
		messages <- message

		Convey("When consume is called", func() {

			consumer := event.NewConsumer()
			consumer.Consume(messageConsumer, handler)

			waitForEventsToBeSentToHandler(handler)

			Convey("A event is sent to the handler ", func() {
				So(len(handler.Events), ShouldEqual, 1)

				event := handler.Events[0]
				So(event.FileURL, ShouldEqual, expectedEvent.FileURL)
				So(event.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})

			Convey("The message is committed", func() {
				So(message.Committed(), ShouldEqual, true)
			})
		})
	})
}

func TestToEvent(t *testing.T) {

	Convey("Given a event schema encoded using avro", t, func() {

		expectedEvent := getExampleEvent()
		message := kafkatest.NewMessage(marshal(*expectedEvent))

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

	Convey("Given a consumer", t, func() {

		messages := make(chan kafka.Message, 1)
		messageConsumer := kafkatest.NewMessageConsumer(messages)
		handler := eventtest.NewEventHandler()

		expectedEvent := getExampleEvent()

		message := kafkatest.NewMessage(marshal(*expectedEvent))
		messages <- message

		consumer := event.NewConsumer()
		consumer.Consume(messageConsumer, handler)

		Convey("When close is called", func() {

			err := consumer.Close(context.Background())

			Convey("The expected event is sent to the handler", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

// marshal helper method to marshal a event into a []byte
func marshal(event event.DimensionsInserted) []byte {
	bytes, err := schema.DimensionsInsertedEvent.Marshal(event)
	So(err, ShouldBeNil)
	return bytes
}

func getExampleEvent() *event.DimensionsInserted {
	expectedEvent := &event.DimensionsInserted{
		InstanceID: "1234",
		FileURL:    "s3://some-file",
	}
	return expectedEvent
}

func waitForEventsToBeSentToHandler(eventHandler *eventtest.EventHandler) {

	start := time.Now()
	timeout := start.Add(time.Millisecond * 500)
	for {
		if len(eventHandler.Events) > 0 {
			log.Debug("events have been sent to the handler", nil)
			break
		}

		if time.Now().After(timeout) {
			log.Debug("timeout hit", nil)
			break
		}

		time.Sleep(time.Millisecond * 10)
	}
}
