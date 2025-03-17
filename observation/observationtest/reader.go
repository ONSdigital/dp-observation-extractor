package observationtest

import (
	"io"

	"github.com/ONSdigital/dp-observation-extractor/observation"
)

var _ observation.Reader = (*Reader)(nil)

// Reader is a reader that returns the given observations and error on read.
type Reader struct {
	offset       int
	observations []*observation.Observation
	error
}

// NewReader provides a reader that returns the given observations and error on read.
func NewReader(observations []*observation.Observation, err error) *Reader {
	return &Reader{
		observations: observations,
		offset:       0,
		error:        err,
	}
}

// Read an observation from the mocked observations.
func (reader *Reader) Read() (*observation.Observation, error) {
	if reader.offset == len(reader.observations) {
		return nil, io.EOF
	}

	observation := reader.observations[reader.offset]
	reader.offset++
	return observation, reader.error
}
