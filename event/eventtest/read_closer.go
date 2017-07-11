package eventtest

import "io"

var _ io.ReadCloser = (*ReadCloser)(nil)

// ReadCloser provides a mock reader that returns the provided error when calling read, and captures a call to close.
type ReadCloser struct {
	ReadError  error
	CloseError error
	Closed     bool
	io.Reader
}

// Read returns the configured error.
func (reader *ReadCloser) Read(p []byte) (n int, err error) {
	if reader.Reader != nil {
		return reader.Reader.Read(p)
	}

	return 0, reader.ReadError
}

// Close captures the fact that it has been called for assertion
func (reader *ReadCloser) Close() (err error) {
	reader.Closed = true
	return reader.CloseError
}
