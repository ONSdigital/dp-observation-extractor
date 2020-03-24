package event

import (
	"context"
	"encoding/hex"
	"io"
	"strconv"

	"github.com/ONSdigital/dp-observation-extractor/observation"
	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go/aws/session"
)

//go:generate moq -out mocks/s3client.go -pkg mock . S3Client
//go:generate moq -out mocks/vault.go -pkg mock . VaultClient

// CSVHandler handles events to extract observations from CSV files.
type CSVHandler struct {
	AwsSession        *session.Session
	s3Clients         map[string]S3Client
	vaultClient       VaultClient
	vaultPath         string
	observationWriter ObservationWriter
}

// NewCSVHandler returns a new CSVHandler instance that uses the given file.FileGetter and Output producer.
func NewCSVHandler(awsSession *session.Session, s3Clients map[string]S3Client, vaultClient VaultClient, observationWriter ObservationWriter, vaultPath string) *CSVHandler {
	return &CSVHandler{
		AwsSession:        awsSession,
		s3Clients:         s3Clients,
		vaultClient:       vaultClient,
		vaultPath:         vaultPath,
		observationWriter: observationWriter,
	}
}

// S3Client represents the S3 client from dp-s3 with the required methods
type S3Client interface {
	Get(key string) (io.ReadCloser, *int64, error)
	GetWithPSK(key string, psk []byte) (io.ReadCloser, *int64, error)
}

// VaultClient is an interface to represent methods called to action upon vault
type VaultClient interface {
	ReadKey(path, key string) (string, error)
}

// ObservationWriter provides operations for observation output.
type ObservationWriter interface {
	WriteAll(ctx context.Context, observationReader observation.Reader, instanceID string)
}

// Handle takes a single event, and returns the observations gathered from the URL in the event.
func (handler CSVHandler) Handle(ctx context.Context, event *DimensionsInserted) error {
	url := event.FileURL

	logData := log.Data{"url": url, "event": event}
	log.Event(ctx, "getting file", log.INFO, logData)

	// parse the url - expected format; s3://bucket/k/e/y
	s3Url, err := s3client.ParseURL(url, s3client.AliasVirtualHostedStyle)
	if err != nil {
		log.Event(ctx, "unable to find bucket and filename in event file url", log.ERROR, log.Error(err), logData)
		return err
	}
	logData["bucket"] = s3Url.BucketName
	logData["filename"] = s3Url.Key

	// Get S3 Client corresponding to the Bucket extracted from URL, or create one if not available
	s3, ok := handler.s3Clients[s3Url.BucketName]
	if !ok {
		log.Event(ctx, "retreiving data from unexpected s3 bucket", log.WARN, log.Data{"RequestedBucket": s3Url.BucketName})
		s3 = s3client.NewClientWithSession(s3Url.BucketName, handler.vaultClient != nil, handler.AwsSession)
	}

	var file io.ReadCloser
	var contentLength *int64
	if handler.vaultClient != nil {

		vaultPath := handler.vaultPath + "/" + s3Url.Key
		vaultKey := "key"
		logData["vault_path"] = vaultPath

		log.Event(ctx, "attempting to get psk from vault", log.INFO, logData)
		pskStr, err := handler.vaultClient.ReadKey(vaultPath, vaultKey)
		if err != nil {
			return err
		}

		log.Event(ctx, "got psk", log.INFO, logData)
		psk, err := hex.DecodeString(pskStr)
		if err != nil {
			return err
		}

		log.Event(ctx, "attempting to get S3 object with psk", log.INFO, logData)

		file, contentLength, err = s3.GetWithPSK(s3Url.Key, psk)
		if err != nil {
			log.Event(ctx, "encountered error retrieving and decrypting csv file", log.ERROR, log.Error(err), logData)
			return err
		}
	} else {
		log.Event(ctx, "attempting to get S3 object", log.INFO, logData)
		file, contentLength, err = s3.Get(s3Url.Key)
		if err != nil {
			log.Event(ctx, "unable to retrieve s3 output object", log.ERROR, log.Error(err), logData)
			return err
		}
	}
	defer file.Close()

	logData["content_length"] = getContentLengthStr(contentLength)
	log.Event(ctx, "file read from s3", log.INFO, logData)

	observationReader := observation.NewCSVReader(file)

	handler.observationWriter.WriteAll(ctx, observationReader, event.InstanceID)
	return nil
}

// getContentLengthStr returns the string representation of the provided *int64, returning '0' if it is nil
func getContentLengthStr(cLen *int64) string {
	if cLen == nil {
		return "0"
	}
	return strconv.FormatInt(*cLen, 10)
}
