package event_test

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/event/eventtest"
	mock "github.com/ONSdigital/dp-observation-extractor/event/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

// S3 testing vars
var (
	bucket          = "some-bucket"
	filename        = "some-file"
	exampleHeader   = "Observation,other,stuff"
	exampleCsvLine  = "153223,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,All categories: Residence Type,All categories: Residence Type,,,"
	contentLen      = int64(284)
	cryptoClientErr = errors.New("crypto client error")
)

// Vault testing vars
var (
	psk        = []byte("Hello World")
	encodedPSK = "48656C6C6F20576F726C64"
	invalidPSK = "this-is-not-hex-encoded"
	vaultPath  = "test-path"
	vaultErr   = errors.New("vault client error")
)

// S3 Get function for a successful case
var funcGetValid = func(key string) (io.ReadCloser, *int64, error) {
	return ioutil.NopCloser(strings.NewReader(exampleHeader + "\n" + exampleCsvLine)), &contentLen, nil
}

// S3 Get function for an error case
var funcGetErr = func(key string) (io.ReadCloser, *int64, error) {
	return nil, nil, io.EOF
}

// S3 GetWithPsk for a successful case
var funcGetWithPskValid = func(key string, psk []byte) (io.ReadCloser, *int64, error) {
	return ioutil.NopCloser(strings.NewReader(exampleHeader + "\n" + exampleCsvLine)), &contentLen, nil
}

// S3 GetWithPsk for an error case
var funcGetWithPskErr = func(key string, psk []byte) (io.ReadCloser, *int64, error) {
	return nil, nil, cryptoClientErr
}

// createS3MockGet creates an S3Client mock with the provided function and returns it, and the map as expected by Handler
func createS3MockGet(funcGet func(key string) (io.ReadCloser, *int64, error)) (*mock.S3ClientMock, map[string]event.S3Client) {
	s3cli := &mock.S3ClientMock{GetFunc: funcGet}
	s3Clients := map[string]event.S3Client{bucket: s3cli}
	return s3cli, s3Clients
}

func createS3MockGetWithPsk(funcGetWithPsk func(key string, psk []byte) (io.ReadCloser, *int64, error)) (*mock.S3ClientMock, map[string]event.S3Client) {
	s3cli := &mock.S3ClientMock{GetWithPSKFunc: funcGetWithPsk}
	s3Clients := map[string]event.S3Client{bucket: s3cli}
	return s3cli, s3Clients
}

// createS3MockEmpty returns an empty s3 mock and the map as expected by Handler
func createS3MockEmpty() (*mock.S3ClientMock, map[string]event.S3Client) {
	s3cli := &mock.S3ClientMock{}
	s3Clients := map[string]event.S3Client{bucket: s3cli}
	return s3cli, s3Clients
}

// Vault ReadKey function for a successful case
var funcReadKey = func(path string, key string) (string, error) {
	return encodedPSK, nil
}

// Vault ReadKey function for an error case
var funcReadKeyErr = func(path string, key string) (string, error) {
	return "", vaultErr
}

// Vault ReadKey function which returns an invalid psk
var funcReadKeyInvalidPSK = func(path string, key string) (string, error) {
	return invalidPSK, nil
}

// createVaultMock creates a Vault mock with the provided function
func createVaultMock(funcReadKey func(path string, key string) (string, error)) *mock.VaultClientMock {
	return &mock.VaultClientMock{ReadKeyFunc: funcReadKey}
}

func TestSuccessfullyHandleCSV(t *testing.T) {

	Convey("Given a valid event message", t, func() {

		Convey("When handle method is called with event", func() {

			Convey("Then successfully return without an error", func() {
				s3cli, s3Clients := createS3MockGet(funcGetValid)
				observationWriterStub := &eventtest.ObservationWriter{}
				csvHandler := event.NewCSVHandler(nil, s3Clients, nil, observationWriterStub, "")

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldBeNil)
				So(len(s3cli.GetCalls()), ShouldEqual, 1)
				So(s3cli.GetCalls()[0].Key, ShouldEqual, filename)
			})
		})

		Convey("When handle method is called with event, and encryption is enabled", func() {

			Convey("Then successfully return without an error", func() {
				s3cli, s3Clients := createS3MockGetWithPsk(funcGetWithPskValid)
				vaultClient := &mock.VaultClientMock{
					ReadKeyFunc: func(path string, key string) (string, error) {
						return encodedPSK, nil
					},
				}

				observationWriterStub := &eventtest.ObservationWriter{}
				csvHandler := event.NewCSVHandler(nil, s3Clients, vaultClient, observationWriterStub, vaultPath)

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldBeNil)

				So(len(s3cli.GetWithPSKCalls()), ShouldEqual, 1)
				So(s3cli.GetWithPSKCalls()[0].Key, ShouldEqual, filename)
				So(s3cli.GetWithPSKCalls()[0].Psk, ShouldResemble, psk)

				So(len(vaultClient.ReadKeyCalls()), ShouldEqual, 1)
				So(vaultClient.ReadKeyCalls()[0].Key, ShouldEqual, "key")
				So(vaultClient.ReadKeyCalls()[0].Path, ShouldEqual, vaultPath+"/"+filename)
			})
		})
	})
}

func TestFailToHandleCSV(t *testing.T) {

	t.Parallel()
	Convey("Given an event is missing a file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'could not find bucket or filename in file url'", func() {
				_, s3Clients := createS3MockEmpty()
				csvHandler := event.NewCSVHandler(nil, s3Clients, nil, &eventtest.ObservationWriter{}, "")
				err := csvHandler.Handle(ctx, &event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "",
				})
				So(err, ShouldResemble, errors.New("wrong bucketName in DNS-alias-virtual-hosted-style url: "))
			})
		})
	})

	Convey("Given an event is missing the filename in file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'missing filename in file url'", func() {
				_, s3Clients := createS3MockEmpty()
				csvHandler := event.NewCSVHandler(nil, s3Clients, nil, &eventtest.ObservationWriter{}, "")
				err := csvHandler.Handle(ctx, &event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "s3://",
				})
				So(err, ShouldResemble, errors.New("wrong bucketName in DNS-alias-virtual-hosted-style url: s3://"))
			})
		})
	})

	Convey("Given a handler with an erroring vault client", t, func() {
		Convey("When handle method is called with a valid event", func() {
			Convey("Then the correct error is returned", func() {
				_, s3Clients := createS3MockEmpty()
				vaultClient := createVaultMock(funcReadKeyErr)
				csvHandler := event.NewCSVHandler(nil, s3Clients, vaultClient, nil, vaultPath)

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldResemble, vaultErr)

				So(len(vaultClient.ReadKeyCalls()), ShouldEqual, 1)
				So(vaultClient.ReadKeyCalls()[0].Key, ShouldEqual, "key")
				So(vaultClient.ReadKeyCalls()[0].Path, ShouldEqual, vaultPath+"/"+filename)
			})
		})
	})

	Convey("Given a handler with an vault client that returns an invalid psk", t, func() {
		Convey("When handle method is called with a valid event", func() {
			Convey("Then the correct error is returned", func() {
				_, s3Clients := createS3MockEmpty()
				vaultClient := createVaultMock(funcReadKeyInvalidPSK)
				csvHandler := event.NewCSVHandler(nil, s3Clients, vaultClient, nil, vaultPath)

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "encoding/hex: invalid byte: U+0074 't'")

				So(len(vaultClient.ReadKeyCalls()), ShouldEqual, 1)
				So(vaultClient.ReadKeyCalls()[0].Key, ShouldEqual, "key")
				So(vaultClient.ReadKeyCalls()[0].Path, ShouldEqual, vaultPath+"/"+filename)
			})
		})
	})

	Convey("Given a handler with a crypto client that returns an error", t, func() {
		Convey("When handle method is called with a valid event", func() {

			Convey("Then the correct error is returned", func() {
				s3cli, s3Clients := createS3MockGetWithPsk(funcGetWithPskErr)
				vaultClient := createVaultMock(funcReadKey)
				csvHandler := event.NewCSVHandler(nil, s3Clients, vaultClient, &eventtest.ObservationWriter{}, vaultPath)

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldResemble, cryptoClientErr)

				So(len(s3cli.GetWithPSKCalls()), ShouldEqual, 1)
				So(s3cli.GetWithPSKCalls()[0].Key, ShouldEqual, filename)
				So(s3cli.GetWithPSKCalls()[0].Psk, ShouldResemble, psk)

				So(len(vaultClient.ReadKeyCalls()), ShouldEqual, 1)
				So(vaultClient.ReadKeyCalls()[0].Key, ShouldEqual, "key")
				So(vaultClient.ReadKeyCalls()[0].Path, ShouldEqual, vaultPath+"/"+filename)
			})
		})
	})

	Convey("Given an event is missing the bucket name in file URL", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'wrong key in global virtual hosted style url: s3://some-file'", func() {
				_, s3Clients := createS3MockEmpty()
				csvHandler := event.NewCSVHandler(nil, s3Clients, nil, &eventtest.ObservationWriter{}, "")

				err := csvHandler.Handle(ctx, &event.DimensionsInserted{
					InstanceID: "1234",
					FileURL:    "s3://some-file",
				})
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("wrong key in global virtual hosted style url: s3://some-file"))
			})
		})
	})

	Convey("Given the event message contains an invalid file url", t, func() {
		Convey("When handle method is called with event", func() {
			Convey("Then error returns with a message 'EOF'", func() {
				s3cli, s3Clients := createS3MockGet(funcGetErr)
				csvHandler := event.NewCSVHandler(nil, s3Clients, nil, &eventtest.ObservationWriter{}, "")

				err := csvHandler.Handle(ctx, getExampleEvent())
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("EOF"))

				So(len(s3cli.GetCalls()), ShouldEqual, 1)
				So(s3cli.GetCalls()[0].Key, ShouldEqual, filename)
			})
		})
	})
}
