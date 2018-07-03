package event

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/ONSdigital/dp-observation-extractor/observation"
	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CSVHandler handles events to extract observations from CSV files.
type CSVHandler struct {
	client            Client
	vaultClient       VaultClient
	observationWriter ObservationWriter

	vaultPath string
}

// NewCSVHandler returns a new CSVHandler instance that uses the given file.FileGetter and Output producer.
func NewCSVHandler(sdkCli SDKClient, cryptoCli CryptoClient, vaultClient VaultClient, observationWriter ObservationWriter, vaultPath string) *CSVHandler {
	return &CSVHandler{
		observationWriter: observationWriter,
		client:            client{SDKClient: sdkCli, CryptoClient: cryptoCli},
		vaultClient:       vaultClient,
		vaultPath:         vaultPath,
	}
}

// Client implements SDKClient and CryptoClient.
type Client interface {
	SDKClient
	CryptoClient
}

// VaultClient is an interface to represent methods called to action upon vault
type VaultClient interface {
	ReadKey(path, key string) (string, error)
}

// client implements the Client interface.
type client struct {
	SDKClient
	CryptoClient
}

// SDKClient retrieves an object from s3.
type SDKClient interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// CryptoClient retrieves an object from s3 that is encrypted.
type CryptoClient interface {
	GetObjectWithPSK(input *s3.GetObjectInput, psk []byte) (*s3.GetObjectOutput, error)
}

// ObservationWriter provides operations for observation output.
type ObservationWriter interface {
	WriteAll(observationReader observation.Reader, instanceID string)
}

// Handle takes a single event, and returns the observations gathered from the URL in the event.
func (handler CSVHandler) Handle(event *DimensionsInserted) error {
	url := event.FileURL

	logData := log.Data{"url": url, "event": event}
	log.Debug("getting file", logData)

	bucket, filename, err := GetBucketAndFilename(url)
	if err != nil {
		log.ErrorC("unable to find bucket and filename in event file url", err, logData)
		return err
	}

	logData["bucket"] = bucket
	logData["filename"] = filename

	getInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filename,
	}

	var output *s3.GetObjectOutput
	if handler.vaultClient != nil {

		vaultPath := handler.vaultPath + "/" + filename
		vaultKey := "key"

		pskStr, err := handler.vaultClient.ReadKey(vaultPath, vaultKey)
		if err != nil {
			return err
		}
		psk, err := hex.DecodeString(pskStr)
		if err != nil {
			return err
		}

		output, err = handler.client.GetObjectWithPSK(getInput, psk)
		if err != nil {
			log.ErrorC("encountered error retrieving and decrypting csv file", err, logData)
			return err
		}
	} else {
		output, err = handler.client.GetObject(getInput)
		if err != nil {
			log.ErrorC("unable to retrieve s3 output object", err, logData)
			return err
		}
	}

	file := output.Body
	defer output.Body.Close()

	logData["content_length"] = strconv.FormatInt(*output.ContentLength, 10)
	log.Info("file read from s3", logData)

	observationReader := observation.NewCSVReader(file)

	handler.observationWriter.WriteAll(observationReader, event.InstanceID)
	return nil
}

// GetBucketAndFilename finds the bucket and filename in url
// FIXME GetBucketAndFilename will fail to retrieve correct file location if folder
// structure is to be introduced in s3 bucket
func GetBucketAndFilename(s3URL string) (string, string, error) {
	urlSplitz := strings.Split(s3URL, "/")
	n := len(urlSplitz)
	if n < 3 {
		return "", "", errors.New("could not find bucket or filename in file url")
	}
	bucket := urlSplitz[n-2]
	filename := urlSplitz[n-1]
	if filename == "" {
		return "", "", errors.New("missing filename in file url")
	}
	if bucket == "" {
		return "", "", errors.New("missing bucket name in file url")
	}

	return bucket, filename, nil
}
