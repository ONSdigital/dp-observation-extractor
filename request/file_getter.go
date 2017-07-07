package request

import (
	"io"
	"github.com/ONSdigital/go-ns/s3"
)

var _ FileGetter = (*s3.S3)(nil)

// FileGetter gets a file reader for the given URL.
type FileGetter interface {
	Get(url string) (io.ReadCloser, error)
}
