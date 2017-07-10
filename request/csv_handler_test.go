package request_test

import (
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/dp-observation-extractor/request/requesttest"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"strings"
	"testing"
)

var exampleHeader string = "Observation,other,stuff"
var exampleCsvLine string = "153223,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,All categories: Residence Type,All categories: Residence Type,,,"

func TestCsvHandler_FileGetterError(t *testing.T) {

	Convey("Given an example request and a file getter that returns an error", t, func() {

		exampleRequest := getExampleRequest()

		expectedError := io.EOF
		fileGetterStub := &requesttest.FileGetter{Error: expectedError}
		observationWriterStub := &requesttest.ObservationWriter{}

		requestHandler := request.NewCSVHandler(fileGetterStub, observationWriterStub)

		Convey("When handle is called", func() {

			_ = requestHandler.Handle(exampleRequest)

			Convey("The error returned from thefile getter is returned from the handler", func() {
				So(fileGetterStub.Error, ShouldEqual, expectedError)
			})
		})
	})
}

func TestCsvHandler_FileGetterUrl(t *testing.T) {

	Convey("Given an example request", t, func() {

		exampleRequest := getExampleRequest()

		fileGetterStub := &requesttest.FileGetter{Error: io.EOF}
		observationWriterStub := &requesttest.ObservationWriter{}

		requestHandler := request.NewCSVHandler(fileGetterStub, observationWriterStub)

		Convey("When handle is called with the example request", func() {

			_ = requestHandler.Handle(exampleRequest)

			Convey("The file getter is called with the url from the request", func() {
				So(fileGetterStub.URL, ShouldEqual, exampleRequest.FileURL)
			})
		})
	})
}

func TestCsvHandler_FileReaderClose(t *testing.T) {

	Convey("Given an example request", t, func() {

		exampleRequest := getExampleRequest()

		stubReadCloser := &requesttest.ReadCloser{}
		fileGetterStub := &requesttest.FileGetter{Reader: stubReadCloser}
		observationWriterStub := &requesttest.ObservationWriter{}

		requestHandler := request.NewCSVHandler(fileGetterStub, observationWriterStub)

		Convey("When handle is called", func() {

			_ = requestHandler.Handle(exampleRequest)

			Convey("The file reader that is returned from the file getter is closed by the handler", func() {
				So(stubReadCloser.Closed, ShouldBeTrue)
			})
		})
	})
}

func TestCsvHandler(t *testing.T) {

	Convey("Given an example request", t, func() {

		exampleRequest := getExampleRequest()

		stubReadCloser := &requesttest.ReadCloser{Reader: strings.NewReader(exampleHeader + "\n" + exampleCsvLine)}
		fileGetterStub := &requesttest.FileGetter{Reader: stubReadCloser}
		observationWriterStub := &requesttest.ObservationWriter{}

		requestHandler := request.NewCSVHandler(fileGetterStub, observationWriterStub)

		Convey("When handle is called", func() {

			err := requestHandler.Handle(exampleRequest)
			So(err, ShouldBeNil)

			Convey("The observation read from the CSV file is sent to the observation writer", func() {

				// first observation is the second row of the CSV (header row discarded)
				observation, err := observationWriterStub.Reader.Read()
				So(err, ShouldBeNil)
				So(observation.Row, ShouldEqual, exampleCsvLine)

				// A second read returns no more observations and EOF error.
				observation2, err := observationWriterStub.Reader.Read()
				So(observation2, ShouldBeNil)
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}
