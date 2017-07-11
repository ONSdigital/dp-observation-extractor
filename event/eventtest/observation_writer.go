package eventtest

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/event"
)

var _ event.ObservationWriter = (*ObservationWriter)(nil)

// ObservationWriter when used will capture the reader passed to it for assertions. Will return the configured error.
type ObservationWriter struct {
	Reader observation.Reader
}

// WriteAll will capture the reader passed to it for assertions.
func (observationWriter *ObservationWriter) WriteAll(reader observation.Reader, instanceID string) {
	observationWriter.Reader = reader
}
