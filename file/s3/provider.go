package s3

import (
	"github.com/ONSdigital/dp-observation-extractor/file"
	"io"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

// NewProvider returns a new AWS specific file.Provider instance configured for the given region.
func NewProvider(awsRegion string) (file.Provider, error) {

	auth, err := aws.SharedAuth()
	if err != nil {
		return nil, err
	}

	s3 := s3.New(auth, aws.Regions[awsRegion])

	return &provider{
		s3: s3,
	}, nil
}

type provider struct {
	s3 *s3.S3
}

// Get a io.ReadCloser instance for the given fully qualified S3 URL
func (provider *provider) Get(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewS3URL(rawURL)
	if err != nil {
		return nil, err
	}

	bucket := provider.s3.Bucket(url.BucketName())
	reader, error := bucket.GetReader(url.Path())
	if err != nil {
		return nil, error
	}

	return reader, nil
}
