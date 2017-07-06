package file

import "io"

// Provider is a generic interface for any provider of files.
type Provider interface {
	Get(url string) (io.ReadCloser, error)
}
