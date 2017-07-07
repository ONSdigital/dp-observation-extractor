package message_test

import (
	"github.com/ONSdigital/dp-observation-extractor/message"
	"github.com/ONSdigital/dp-observation-extractor/mock"
	"github.com/ONSdigital/dp-observation-extractor/model"
	"github.com/ONSdigital/go-ns/avro"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/ONSdigital/dp-observation-extractor/message/messagetest"
)

func TestConsumeMessages_UnmarshallError(t *testing.T) {
	Convey("Given a message consumer with an invalid message and a valid message", t, func() {

		messages := make(chan []byte, 2)
		messageConsumer := messagetest.NewMessageConsumer(messages)
		requestHandler := mock.NewRequestHandler()

		expectedRequest := model.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		messages <- []byte("invalid message")
		messages <- bytes
		close(messages)

		Convey("When consume messages is called", func() {

			message.ConsumeMessages(messageConsumer, requestHandler)

			Convey("Only the valid request is sent to the handler ", func() {
				So(len(requestHandler.Requests), ShouldEqual, 1)

				request := requestHandler.Requests[0]
				So(request.FileURL, ShouldEqual, expectedRequest.FileURL)
				So(request.InstanceID, ShouldEqual, expectedRequest.InstanceID)
			})
		})
	})
}

func TestConsumeMessages(t *testing.T) {

	Convey("Given a message consumer with a valid message", t, func() {

		messages := make(chan []byte, 1)
		messageConsumer := messagetest.NewMessageConsumer(messages)
		requestHandler := mock.NewRequestHandler()

		expectedRequest := model.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		messages <- bytes
		close(messages)

		Convey("When consume messages is called", func() {

			message.ConsumeMessages(messageConsumer, requestHandler)

			Convey("A request is sent to the handler ", func() {
				So(len(requestHandler.Requests), ShouldEqual, 1)

				request := requestHandler.Requests[0]
				So(request.FileURL, ShouldEqual, expectedRequest.FileURL)
				So(request.InstanceID, ShouldEqual, expectedRequest.InstanceID)
			})
		})
	})
}

func TestToRequest(t *testing.T) {

	Convey("Given a request message encoded using avro", t, func() {

		expectedRequest := model.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		Convey("When the expectedRequest is unmarshalled", func() {

			request, err := message.ToRequest(bytes)

			Convey("The expectedRequest has the expected values", func() {
				So(err, ShouldBeNil)
				So(request.FileURL, ShouldEqual, expectedRequest.FileURL)
				So(request.InstanceID, ShouldEqual, expectedRequest.InstanceID)
			})
		})
	})
}

// Helper method to marshal a request into a []byte
func toBytes(request model.Request) ([]byte, error) {
	marshalSchema := &avro.Schema{
		Definition: message.RequestSchema,
	}
	bytes, err := marshalSchema.Marshal(request)
	return bytes, err
}
