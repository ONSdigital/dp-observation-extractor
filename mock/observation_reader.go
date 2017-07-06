package mock

import (
	"github.com/ONSdigital/dp-observation-extractor/model"
	"github.com/ONSdigital/dp-observation-extractor/observation"
	"io"
)

var _ observation.Reader = (*observationReader)(nil)

func NewObservationReader(observations []*model.Observation, error error) observation.Reader {
	return &observationReader{
		observations: observations,
		offset:       0,
		error:        error,
	}
}

type observationReader struct {
	offset       int
	observations []*model.Observation
	error
}

func (reader *observationReader) Read() (*model.Observation, error) {
	if reader.offset == len(reader.observations) {
		return nil, io.EOF
	}

	observation := reader.observations[reader.offset]
	reader.offset++
	return observation, reader.error
}
