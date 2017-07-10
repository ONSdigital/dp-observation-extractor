package request_test

import (
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/dp-observation-extractor/request/requesttest"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConsumeMessages_UnmarshallError(t *testing.T) {
	Convey("Given a schema consumer with an invalid schema and a valid schema", t, func() {

		messages := make(chan []byte, 2)
		messageConsumer := requesttest.NewMessageConsumer(messages)
		requestHandler := requesttest.NewRequestHandler()

		expectedRequest := getExampleRequest()

		messages <- []byte("invalid schema")
		messages <- Marshal(*expectedRequest)
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

	Convey("Given a schema consumer with a valid schema", t, func() {

		messages := make(chan []byte, 1)
		messageConsumer := requesttest.NewMessageConsumer(messages)
		requestHandler := requesttest.NewRequestHandler()

		expectedRequest := getExampleRequest()

		messages <- Marshal(*expectedRequest)
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

	Convey("Given a request schema encoded using avro", t, func() {

		expectedRequest := getExampleRequest()
		bytes := Marshal(*expectedRequest)

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

// Marshal helper method to marshal a request into a []byte
func Marshal(request request.Request) []byte {
	bytes, err := schema.Request.Marshal(request)
	So(err, ShouldBeNil)
	return bytes
}

func getExampleRequest() *request.Request {
	expectedRequest := &request.Request{
		InstanceID: "1234",
		FileURL:    "s3://some-file",
	}
	return expectedRequest
}
