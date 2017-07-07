package observation


// BatchReader provides an interface to read Observations
type BatchReader interface {
	Read(batchSize int) ([]*Observation, error)
}

// ensure CSVReader satisfies the CSVReader interface
var _ BatchReader = (*batchReader)(nil)

// BatchReader reads the specified number of observation
type batchReader struct {
	reader Reader
}

// NewBatchReader returns a new BatchReader instance for the given CSVReader
func NewBatchReader(observationReader Reader) BatchReader {
	return &batchReader{
		reader: observationReader,
	}
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (batchReader *batchReader) Read(batchSize int) ([]*Observation, error) {

	batch := make([]*Observation, 0, batchSize)

	for index := 0; index < batchSize; index++ {

		observation, err := batchReader.reader.Read()
		if err != nil {
			return batch, err
		}

		batch = append(batch, observation)
	}

	return batch, nil
}
