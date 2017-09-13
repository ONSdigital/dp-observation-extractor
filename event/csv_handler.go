package event

import (
	"io"

	"github.com/ONSdigital/dp-observation-extractor/errors"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/go-ns/log"
)

// CSVHandler handles events to extract observations from CSV files.
type CSVHandler struct {
	fileGetter        FileGetter
	observationWriter ObservationWriter
	errorHandler      errors.Handler
}

// NewCSVHandler returns a new CSVHandler instance that uses the given file.FileGetter and Output producer.
func NewCSVHandler(fileGetter FileGetter, observationWriter ObservationWriter, errorHandler errors.Handler) *CSVHandler {
	return &CSVHandler{
		observationWriter: observationWriter,
		fileGetter:        fileGetter,
		errorHandler:      errorHandler,
	}
}

// FileGetter gets a file reader for the given URL.
type FileGetter interface {
	Get(url string) (io.ReadCloser, error)
}

// ObservationWriter provides operations for observation output.
type ObservationWriter interface {
	WriteAll(observationReader observation.Reader, instanceID string)
}

// Handle takes a single event, and returns the observations gathered from the URL in the event.
func (handler CSVHandler) Handle(event *DimensionsInserted) error {

	url := event.FileURL

	log.Debug("getting file", log.Data{"url": url, "event": event})
	readCloser, err := handler.fileGetter.Get(url)
	if err != nil {
		handler.errorHandler.Handle(event.InstanceID, err)
		return err
	}
	defer func(readCloser io.ReadCloser) {
		closeErr := readCloser.Close()
		if closeErr != nil {
			handler.errorHandler.Handle(event.InstanceID, err)
			log.Error(closeErr, nil)
		}
	}(readCloser)

	observationReader := observation.NewCSVReader(readCloser)

	handler.observationWriter.WriteAll(observationReader, event.InstanceID)
	return nil
}
