package requesttest

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

var _ request.ObservationWriter = (*ObservationWriter)(nil)

// ObservationWriter when used will capture the reader passed to it for assertions. Will return the configured error.
type ObservationWriter struct {
	Reader observation.Reader
	Error  error
}

// WriteAll will capture the reader passed to it for assertions.
func (observationWriter *ObservationWriter) WriteAll(reader observation.Reader) error {
	observationWriter.Reader = reader
	return observationWriter.Error
}
