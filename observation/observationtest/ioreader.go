package observationtest

import "io"

var _ io.Reader = (*IOReader)(nil)

// IOReader provides a mock reader that returns the provided error when calling read.
type IOReader struct {
	err error
}

// NewIOReader returns a new instance of CSVReader for the given error.
func NewIOReader(err error) *IOReader {
	return &IOReader{
		err: err,
	}
}

// Read returns the configured error.
func (reader *IOReader) Read(p []byte) (n int, err error) {
	return 0, reader.err
}
