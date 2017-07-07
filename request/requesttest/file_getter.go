package requesttest

import (
	"io"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

// Check that Provider implements the file.Provider interface.
var _ request.FileGetter = (*Provider)(nil)

// Provider is a mock file provider that returns the stored io.ReadCloser and error on Get().
type Provider struct {
	Reader io.ReadCloser
	Error  error
}

// Get the configured io.ReadCloser and error.
func (provider *Provider) Get(url string) (io.ReadCloser, error) {
	return provider.Reader, provider.Error
}
