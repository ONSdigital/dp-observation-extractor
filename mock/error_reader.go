package mock

// ErrReader provides a mock reader that returns the provided error when calling read.
type ErrReader struct {
	err error
}

// NewErrReader returns a new instance of ErrReader for the given error.
func NewErrReader(err error) *ErrReader {
	return &ErrReader{
		err: err,
	}
}

// Read returns the configured error.
func (reader *ErrReader) Read(p []byte) (n int, err error) {
	return 0, reader.err
}
