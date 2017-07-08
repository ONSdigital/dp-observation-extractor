package requesttest

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

var _ request.ObservationWriter = (*MessageWriter)(nil)

// MessageWriter when used will capture the reader passed to it for assertions. Will return the configured error.
type MessageWriter struct {
	Reader observation.Reader
	Error  error
}

// WriteAll will capture the reader passed to it for assertions.
func (messageWriter *MessageWriter) WriteAll(reader observation.Reader) error {
	messageWriter.Reader = reader
	return messageWriter.Error
}
