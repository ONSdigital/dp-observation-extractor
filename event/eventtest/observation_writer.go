package eventtest

import (
	"context"

	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/observation"
)

var _ event.ObservationWriter = (*ObservationWriter)(nil)

// ObservationWriter when used will capture the reader passed to it for assertions. Will return the configured error.
type ObservationWriter struct {
	Reader observation.Reader
}

// WriteAll will capture the reader passed to it for assertions.
func (observationWriter *ObservationWriter) WriteAll(ctx context.Context, reader observation.Reader, instanceID string) {
	observationWriter.Reader = reader
}
