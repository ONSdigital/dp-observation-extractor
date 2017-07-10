package observation_test

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/observation/observationtest"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMessageWriter_WriteAll(t *testing.T) {

	Convey("Given an observation reader with a single observation", t, func() {

		// Create mock reader with one expected observation.
		expectedObservation := &observation.Observation{Row: "the,row,content"}
		expectedInstanceID := "123abc"
		expectedEvent := observation.ExtractedEvent{Row: expectedObservation.Row, InstanceID: expectedInstanceID}

		expectedObservations := make([]*observation.Observation, 1)
		expectedObservations[0] = expectedObservation
		mockObservationReader := observationtest.NewReader(expectedObservations, nil)

		// mock schema producer contains the output channel to capture messages sent.
		outputChannel := make(chan []byte, 1)
		mockMessageProducer := observationtest.MessageProducer{OutputChannel: outputChannel}

		observationMessageWriter := observation.NewMessageWriter(mockMessageProducer)

		Convey("When write all is called on the observation schema writer", func() {

			observationMessageWriter.WriteAll(mockObservationReader, expectedInstanceID)

			Convey("The schema producer has the observation on its output channel", func() {

				messageBytes := <-outputChannel
				close(outputChannel)
				observationEvent := Unmarshal(messageBytes)
				So(observationEvent.InstanceID, ShouldEqual, expectedEvent.InstanceID)
			})
		})
	})
}

func TestMessageWriter_Marshal(t *testing.T) {

	Convey("Given an example observation", t, func() {

		expectedObservation := &observation.Observation{Row: "the,row,content"}
		expectedInstanceID := "123abc"
		expectedEvent := observation.ExtractedEvent{Row: expectedObservation.Row, InstanceID: expectedInstanceID}

		Convey("When Marshal is called", func() {

			bytes, err := observation.Marshal(expectedEvent)
			So(err, ShouldBeNil)

			Convey("The observation can be unmarshalled and has the expected values", func() {

				actualEvent := Unmarshal(bytes)
				So(actualEvent.InstanceID, ShouldEqual, expectedEvent.InstanceID)
				So(actualEvent.Row, ShouldEqual, expectedEvent.Row)
			})
		})
	})
}

// Unmarshal converts observation events to []byte.
func Unmarshal(bytes []byte) *observation.ExtractedEvent {
	event := &observation.ExtractedEvent{}
	err := schema.ObservationExtractedEvent.Unmarshal(bytes, event)
	So(err, ShouldBeNil)
	return event
}
