package requesttest

import (
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/dp-observation-extractor/observation"
)

var _ request.ObservationSender = (*OutputSender)(nil)

// OutputSender when used will capture the reader passed to it for assertions. Will return the configured error.
type OutputSender struct {
	Reader observation.Reader
	Error  error
}

// SendObservations will capture the reader passed to it for assertions.
func (outputSender *OutputSender) SendObservations(reader observation.Reader) error {
	outputSender.Reader = reader
	return outputSender.Error
}
