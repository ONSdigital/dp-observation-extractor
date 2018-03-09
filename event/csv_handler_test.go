package event_test

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/event/eventtest"
	"github.com/ONSdigital/dp-observation-extractor/event/mocks"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	bucket         = "some-bucket"
	filename       = "some-file"
	exampleHeader  = "Observation,other,stuff"
	exampleCsvLine = "153223,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,All categories: Residence Type,All categories: Residence Type,,,"
)

var getInput = &s3.GetObjectInput{
	Bucket: &bucket,
	Key:    &filename,
}

func TestSuccessfullyHandleCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a valid event message", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then successfully return without an error", func() {
				client := mocks.NewMockS3API(mockCtrl)
				observationWriterStub := &eventtest.ObservationWriter{}

				s3GetOutput := &s3.GetObjectOutput{
					Body: ioutil.NopCloser(strings.NewReader(exampleHeader + "\n" + exampleCsvLine)),
				}

				client.EXPECT().GetObject(getInput).Return(s3GetOutput, nil)

				csvHandler := event.NewCSVHandler(client, observationWriterStub)

				err := csvHandler.Handle(getExampleEvent())

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestFailToHandleCSV(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Parallel()
	Convey("Given an event is missing a file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'could not find bucket or filename in file url'", func() {
				client := mocks.NewMockS3API(mockCtrl)
				observationWriterStub := &eventtest.ObservationWriter{}

				csvHandler := event.NewCSVHandler(client, observationWriterStub)

				err := csvHandler.Handle(&event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "",
				})

				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("could not find bucket or filename in file url"))
			})
		})
	})

	Convey("Given an event is missing the filename in file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'missing filename in file url'", func() {
				client := mocks.NewMockS3API(mockCtrl)
				observationWriterStub := &eventtest.ObservationWriter{}

				csvHandler := event.NewCSVHandler(client, observationWriterStub)

				err := csvHandler.Handle(&event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "s3://",
				})

				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("missing filename in file url"))
			})
		})
	})

	Convey("Given an event is missing the bucket name in file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'missing bucket name in file url'", func() {
				client := mocks.NewMockS3API(mockCtrl)
				observationWriterStub := &eventtest.ObservationWriter{}

				csvHandler := event.NewCSVHandler(client, observationWriterStub)

				err := csvHandler.Handle(&event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "s3://some-file",
				})

				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("missing bucket name in file url"))
			})
		})
	})

	Convey("Given the event message contains an invalid file url", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'EOF'", func() {
				client := mocks.NewMockS3API(mockCtrl)
				observationWriterStub := &eventtest.ObservationWriter{}
				client.EXPECT().GetObject(getInput).Return(nil, io.EOF)

				csvHandler := event.NewCSVHandler(client, observationWriterStub)

				err := csvHandler.Handle(getExampleEvent())

				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("EOF"))
			})
		})
	})
}

func TestSuccessfullyGetBucketAndFilename(t *testing.T) {
	Convey("Given a valid s3 url (includes bucket and filename)", t, func() {
		Convey("When GetBucketAndFilename function is called", func() {
			Convey("Then function returns both bucket and filename from file url'", func() {
				bucket, filename, err := event.GetBucketAndFilename("s3://some-bucket/some-file")

				So(err, ShouldBeNil)
				So(bucket, ShouldEqual, "some-bucket")
				So(filename, ShouldEqual, "some-file")
			})
		})
	})
}

func TestFailToGetBucketAndFilename(t *testing.T) {

	Convey("Given an empty s3 url", t, func() {
		Convey("When GetBucketAndFilename function is called", func() {
			Convey("Then function responds with an error of 'could not find bucket or filename in file url'", func() {
				bucket, filename, err := event.GetBucketAndFilename("")

				So(bucket, ShouldEqual, "")
				So(filename, ShouldEqual, "")
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("could not find bucket or filename in file url"))
			})
		})
	})

	Convey("Given an s3 url that is missing the bucket name", t, func() {
		Convey("When GetBucketAndFilename function is called", func() {
			Convey("Then function responds with an error of 'missing bucket name in file url'", func() {
				bucket, filename, err := event.GetBucketAndFilename("s3://some-file")

				So(bucket, ShouldEqual, "")
				So(filename, ShouldEqual, "")
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("missing bucket name in file url"))
			})
		})
	})

	Convey("Given an s3 url that is missing the filename", t, func() {
		Convey("When GetBucketAndFilename function is called", func() {
			Convey("Then function responds with an error of 'missing filename in file url'", func() {
				bucket, filename, err := event.GetBucketAndFilename("s3://")

				So(bucket, ShouldEqual, "")
				So(filename, ShouldEqual, "")
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("missing filename in file url"))
			})
		})
	})
}
