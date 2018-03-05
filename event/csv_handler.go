package event

import (
	"errors"
	"strings"

	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CSVHandler handles events to extract observations from CSV files.
type CSVHandler struct {
	client            Client
	observationWriter ObservationWriter
}

// NewCSVHandler returns a new CSVHandler instance that uses the given file.FileGetter and Output producer.
func NewCSVHandler(client Client, observationWriter ObservationWriter) *CSVHandler {
	return &CSVHandler{
		observationWriter: observationWriter,
		client:            client,
	}
}

// Client gets an s3 object.
type Client interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// ObservationWriter provides operations for observation output.
type ObservationWriter interface {
	WriteAll(observationReader observation.Reader, instanceID string)
}

// Handle takes a single event, and returns the observations gathered from the URL in the event.
func (handler CSVHandler) Handle(event *DimensionsInserted) error {
	url := event.FileURL

	logData := log.Data{"url": url, "event": event}
	log.Debug("getting file", logData)

	bucket, filename, err := GetBucketAndFilename(url)
	if err != nil {
		log.ErrorC("unable to find bucket and filename in event file url", err, logData)
		return err
	}

	logData["bucket"] = bucket
	logData["filename"] = filename

	getInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filename,
	}

	output, err := handler.client.GetObject(getInput)
	if err != nil {
		log.ErrorC("unable to retrieve s3 output object", err, logData)
		return err
	}

	file := output.Body
	defer output.Body.Close()

	observationReader := observation.NewCSVReader(file)

	handler.observationWriter.WriteAll(observationReader, event.InstanceID)
	return nil
}

// GetBucketAndFilename finds the bucket and filename in url
func GetBucketAndFilename(s3URL string) (string, string, error) {
	urlSplitz := strings.Split(s3URL, "/")
	n := len(urlSplitz)
	if n < 3 {
		return "", "", errors.New("could not find bucket or filename in file url")
	}
	bucket := urlSplitz[n-2]
	filename := urlSplitz[n-1]
	if filename == "" {
		return "", "", errors.New("missing filename in file url")
	}
	if bucket == "" {
		return "", "", errors.New("missing bucket name in file url")
	}

	return bucket, filename, nil
}
