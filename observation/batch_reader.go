package observation

// BatchReader provides an interface to read Observations
type BatchReader interface {
	Read(batchSize int) ([]*Observation, error)
}

// BatchReader reads the specified number of observation
type batchReader struct {
	observationReader Reader
}

// NewBatchReader returns a new BatchReader instance for the given CSVReader
func NewBatchReader(observationReader Reader) BatchReader {
	return &batchReader{
		observationReader: observationReader,
	}
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (batchReader *batchReader) Read(batchSize int) ([]*Observation, error) {

	batch := make([]*Observation, 0, batchSize)

	for index := 0; index < batchSize; index++ {

		observation, err := batchReader.observationReader.Read()
		if err != nil {
			return batch, err
		}

		batch = append(batch, observation)
	}

	return batch, nil
}
