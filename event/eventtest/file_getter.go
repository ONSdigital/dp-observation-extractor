package eventtest

import (
	"github.com/ONSdigital/dp-observation-extractor/event"
	"io"
)

// Check that FileGetter implements the file.FileGetter interface.
var _ event.FileGetter = (*FileGetter)(nil)

// FileGetter is a mock file getter that returns the stored io.ReadCloser / error on Get(), and captures the last URL passed.
type FileGetter struct {
	URL    string
	Reader io.ReadCloser
	Error  error
}

// Get the configured io.ReadCloser and error.
func (getter *FileGetter) Get(url string) (io.ReadCloser, error) {
	getter.URL = url
	return getter.Reader, getter.Error
}