package requesttest

import (
	"io"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

// Check that Getter implements the file.Getter interface.
var _ request.FileGetter = (*Getter)(nil)

// Getter is a mock file getter that returns the stored io.ReadCloser and error on Get().
type Getter struct {
	Reader io.ReadCloser
	Error  error
}

// Get the configured io.ReadCloser and error.
func (getter *Getter) Get(url string) (io.ReadCloser, error) {
	return getter.Reader, getter.Error
}
