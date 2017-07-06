package observation_test

import (
	"github.com/ONSdigital/dp-observation-extractor/mock"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"strings"
	"testing"
)

var exampleCsvLine string = "153223,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,All categories: Residence Type,All categories: Residence Type,,,"

func TestEmptyInput(t *testing.T) {

	Convey("Given a reader with no content", t, func() {

		reader := strings.NewReader("")
		observationReader := observation.NewReader(reader)

		Convey("When read is called", func() {

			_, err := observationReader.Read()

			Convey("Then an EOF error is returned", func() {
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}

func TestValidInput(t *testing.T) {

	Convey("Given a reader with two rows of data", t, func() {

		reader := strings.NewReader(exampleCsvLine + "\n" + exampleCsvLine)
		observationReader := observation.NewReader(reader)

		Convey("When read is called", func() {

			observation1, err1 := observationReader.Read()
			observation2, err2 := observationReader.Read()
			observation3, err3 := observationReader.Read() // EOF expected

			Convey("There are no errors returned until EOF is reached", func() {

				So(err1, ShouldBeNil)
				So(err2, ShouldBeNil)
				So(err3, ShouldEqual, io.EOF)
			})

			Convey("The two observation instances are populated as expected", func() {

				So(observation1, ShouldNotBeNil)
				So(observation1.Row, ShouldEqual, exampleCsvLine)

				So(observation2, ShouldNotBeNil)
				So(observation2.Row, ShouldEqual, exampleCsvLine)

				So(observation3, ShouldBeNil)
			})
		})
	})
}

func TestErrorResponse(t *testing.T) {

	Convey("Given a reader that returns an error that is not EOF", t, func() {

		expectedError := errors.New("The world has ended")
		observationReader := observation.NewReader(mock.NewErrReader(expectedError))

		Convey("When read is called", func() {

			_, err := observationReader.Read()

			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, expectedError)
			})
		})
	})
}
