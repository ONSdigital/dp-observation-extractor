package s3

import (
	"net/url"
	"strings"
)

// URL represents a fully qualified S3 URL.
type URL interface {
	BucketName() string
	FilePath() string
}

// NewURL creates a new instance of URL.
func NewURL(rawURL string) (URL, error) {

	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &s3URL{
		URL: url,
	}, nil
}

type s3URL struct {
	*url.URL
}

func (url *s3URL) BucketName() string {
	return url.Host
}

func (url *s3URL) FilePath() string {
	return strings.TrimPrefix(url.Path, "/")
}
