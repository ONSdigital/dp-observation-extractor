package request_test

import (
	"github.com/ONSdigital/dp-observation-extractor/message"
	"github.com/ONSdigital/go-ns/avro"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/ONSdigital/dp-observation-extractor/request/requesttest"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

func TestConsumeMessages_UnmarshallError(t *testing.T) {
	Convey("Given a message consumer with an invalid message and a valid message", t, func() {

		messages := make(chan []byte, 2)
		messageConsumer := requesttest.NewMessageConsumer(messages)
		requestHandler := requesttest.NewRequestHandler()

		expectedRequest := request.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		messages <- []byte("invalid message")
		messages <- bytes
		close(messages)

		Convey("When consume messages is called", func() {

			request.Consume(messageConsumer, requestHandler)

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
		messageConsumer := requesttest.NewMessageConsumer(messages)
		requestHandler := requesttest.NewRequestHandler()

		expectedRequest := request.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		messages <- bytes
		close(messages)

		Convey("When consume messages is called", func() {

			request.Consume(messageConsumer, requestHandler)

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

		expectedRequest := request.Request{
			InstanceID: "1234",
			FileURL:    "s3://some-file",
		}

		bytes, err := toBytes(expectedRequest)
		So(err, ShouldBeNil)

		Convey("When the expectedRequest is unmarshalled", func() {

			request, err := request.Unmarshal(bytes)

			Convey("The expectedRequest has the expected values", func() {
				So(err, ShouldBeNil)
				So(request.FileURL, ShouldEqual, expectedRequest.FileURL)
				So(request.InstanceID, ShouldEqual, expectedRequest.InstanceID)
			})
		})
	})
}

// Helper method to marshal a request into a []byte
func toBytes(request request.Request) ([]byte, error) {
	marshalSchema := &avro.Schema{
		Definition: message.RequestSchema,
	}
	bytes, err := marshalSchema.Marshal(request)
	return bytes, err
}
