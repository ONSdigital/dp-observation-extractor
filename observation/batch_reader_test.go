package observation_test

import (
	"github.com/ONSdigital/dp-observation-extractor/model"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"testing"
	"github.com/ONSdigital/dp-observation-extractor/observation/observationtest"
)

func TestBatchReader_Read(t *testing.T) {

	Convey("Given a batch reader with one observation", t, func() {

		expectedObservations := make([]*model.Observation, 1)
		expectedObservations[0] = &model.Observation{Row: "the,row,content"}

		mockObservationReader := observationtest.NewObservationReader(expectedObservations, nil)
		batchReader := observation.NewBatchReader(mockObservationReader)

		Convey("When a batch is read that hits EOF", func() {

			observations, err := batchReader.Read(2)

			Convey("The batch has the observation", func() {
				So(len(observations), ShouldEqual, len(expectedObservations))
				So(observations[0], ShouldEqual, expectedObservations[0])
			})

			Convey("The error value should be EOF", func() {
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}

func TestBatchReader_Read_MultipleBatches(t *testing.T) {

	Convey("Given a batch reader with three observation", t, func() {

		expectedObservations := make([]*model.Observation, 3)
		expectedObservations[0] = &model.Observation{Row: "1,row,content"}
		expectedObservations[1] = &model.Observation{Row: "2,row,content"}
		expectedObservations[2] = &model.Observation{Row: "3,row,content"}

		mockObservationReader := observationtest.NewObservationReader(expectedObservations, nil)
		batchReader := observation.NewBatchReader(mockObservationReader)

		Convey("When multiple batches are read using a batch size of two", func() {

			batchSize := 2
			observations, err := batchReader.Read(batchSize)

			Convey("The first batch has two observation", func() {
				So(len(observations), ShouldEqual, batchSize)
				So(observations[0], ShouldEqual, expectedObservations[0])
				So(observations[1], ShouldEqual, expectedObservations[1])
			})

			Convey("The error value should be nil", func() {
				So(err, ShouldEqual, nil)
			})

			observations, err = batchReader.Read(batchSize)

			Convey("The second batch one observation", func() {
				So(len(observations), ShouldEqual, 1)
				So(observations[0], ShouldEqual, expectedObservations[2])
			})

			Convey("The error value should be EOF", func() {
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}

func TestBatchReader_Read_EOF(t *testing.T) {

	Convey("Given a batch reader with no observation", t, func() {

		mockObservationReader := observationtest.NewObservationReader(make([]*model.Observation, 0), nil)
		batchReader := observation.NewBatchReader(mockObservationReader)

		Convey("When a batch is read", func() {

			_, err := batchReader.Read(1)

			Convey("The error value should be EOF", func() {
				So(err, ShouldEqual, io.EOF)
			})
		})
	})

}
