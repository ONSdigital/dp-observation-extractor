package request

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/go-ns/log"
	"io"
)

// CSVHandler handles requests to extract observations from CSV files.
type CSVHandler struct {
	fileGetter        FileGetter
	observationWriter ObservationWriter
}

// NewCSVHandler returns a new CSVHandler instance that uses the given file.FileGetter and Output producer.
func NewCSVHandler(fileGetter FileGetter, observationWriter ObservationWriter) *CSVHandler {
	return &CSVHandler{
		observationWriter: observationWriter,
		fileGetter:        fileGetter,
	}
}

// FileGetter gets a file reader for the given URL.
type FileGetter interface {
	Get(url string) (io.ReadCloser, error)
}

// ObservationWriter provides operations for observation output.
type ObservationWriter interface {
	WriteAll(observationReader observation.Reader, instanceID string) error
}

// Handle takes a single request, and returns the observations gathered from the URL in the request.
func (handler CSVHandler) Handle(request *Request) error {

	url := request.FileURL

	log.Debug("Getting file", log.Data{"url": url, "request": request})
	readCloser, err := handler.fileGetter.Get(url)
	if err != nil {
		return err
	}
	defer func(readCloser io.ReadCloser) {
		closeErr := readCloser.Close()
		if closeErr != nil {
			log.Error(closeErr, nil)
		}
	}(readCloser)

	discardHeaderRow := true
	observationReader := observation.NewCSVReader(readCloser, discardHeaderRow)

	err = handler.observationWriter.WriteAll(observationReader, request.InstanceID)
	return err
}
