package requesttest

import (
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

var _ request.OutputProducer = (*OutputProducer)(nil)

// OutputProducer when used will capture the reader passed to it for assertions. Will return the configured error.
type OutputProducer struct {
	Reader observation.Reader
	Error error
}

// SendObservations will capture the reader passed to it for assertions.
func (outputProducer *OutputProducer) SendObservations(reader observation.Reader) error {
	outputProducer.Reader = reader
	return outputProducer.Error
}
