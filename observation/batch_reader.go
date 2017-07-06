package observation

import "github.com/ONSdigital/dp-observation-extractor/model"

// BatchReader provides an interface to read Observations
type BatchReader interface {
	Read(batchSize int) ([]*model.Observation, error)
}

// ensure reader satisfies the observation.Reader interface
var _ BatchReader = (*batchReader)(nil)

// BatchReader reads the specified number of observation
type batchReader struct {
	reader Reader
}

// NewBatchReader returns a new BatchReader instance for the given observation.Reader
func NewBatchReader(observationReader Reader) BatchReader {
	return &batchReader{
		reader: observationReader,
	}
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (batchReader *batchReader) Read(batchSize int) ([]*model.Observation, error) {

	batch := make([]*model.Observation, 0, batchSize)

	for index := 0; index < batchSize; index++ {

		observation, err := batchReader.reader.Read()
		if err != nil {
			return batch, err
		}

		batch = append(batch, observation)
	}

	return batch, nil
}
