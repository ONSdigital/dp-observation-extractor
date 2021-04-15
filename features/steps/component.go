package steps

import (
	"context"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/event"
	s3mocks "github.com/ONSdigital/dp-observation-extractor/event/mocks"
	"github.com/ONSdigital/dp-observation-extractor/initialise"
	"github.com/ONSdigital/dp-observation-extractor/initialise/mock"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Component struct {
	componenttest.ErrorFeature
	serviceList   *initialise.ExternalServiceList
	KafkaConsumer kafka.IConsumerGroup
	KafkaProducer kafka.IProducer
	killChannel   chan os.Signal
	apiFeature    *componenttest.APIFeature
	errorChan     chan error
	signals       chan os.Signal
	cfg           *config.Config
	S3Client      *s3mocks.S3ClientMock
}

func NewComponent() *Component {

	c := &Component{errorChan: make(chan error)}

	consumer := kafkatest.NewMessageConsumer(false)
	consumer.CheckerFunc = funcCheck
	c.KafkaConsumer = consumer

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	c.cfg = cfg
	c.S3Client = &s3mocks.S3ClientMock{}

	initialiser := &mock.InitialiserMock{
		DoGetConsumerFunc:      c.DoGetConsumer,
		DoGetHealthCheckerFunc: c.DoGetHealthCheck,
		DoGetProducerFunc:      c.DoGetProducer,
		DoGetS3ClientsFunc:     c.DoGetS3Clients,
	}

	c.serviceList = &initialise.ExternalServiceList{
		Init: initialiser,
	}

	return c
}

func (c *Component) Close() {

}

func (c *Component) Reset() {

}

func (c *Component) DoGetHealthCheck(ctx context.Context, buildTime string, gitCommit string, version string, cfg *config.Config) (*healthcheck.HealthCheck, error) {
	versionInfo, err := healthcheck.NewVersionInfo("1234", "gitCommit", "version")
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCriticalTimeout, cfg.HealthCheckInterval)
	hc.Status = "200"

	return &hc, nil
}

func (c *Component) DoGetConsumer(ctx context.Context, cfg *config.Config) (kafkaConsumer kafka.IConsumerGroup, err error) {
	return c.KafkaConsumer, nil
}

func (c *Component) DoGetProducer(ctx context.Context, topic string, name initialise.KafkaProducerName, cfg *config.Config) (kafkaProducer kafka.IProducer, err error) {
	channels := &kafka.ProducerChannels{
		Output: make(chan []byte),
	}
	c.KafkaProducer = &kafkatest.IProducerMock{
		ChannelsFunc: func() *kafka.ProducerChannels {
			return channels
		},
		CloseFunc:   funcClose,
		CheckerFunc: funcCheck,
	}

	return c.KafkaProducer, nil
}

func (c *Component) DoGetS3Clients(cfg *config.Config) (awsSession *session.Session, s3Clients map[string]event.S3Client, err error) {
	s3Clients = make(map[string]event.S3Client)
	s3Clients["bucket_name"] = c.S3Client
	return
}

func funcClose(ctx context.Context) error {
	return nil
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}

// func NewS3Client() {

// 	return &dps3.S3{
// 		sdkClient:     sdkClient,
// 		cryptoClient:  cryptoClient,
// 		bucketName:    "bucket_name",
// 		region:        region,
// 		mutexUploadID: &sync.Mutex{},
// 		session:       s,
// 	}
// }

// func createS3MockGetWithPsk(funcGetWithPsk func(key string, psk []byte) (io.ReadCloser, *int64, error)) (*mock.S3ClientMock, map[string]event.S3Client) {
// 	s3cli := &mock.S3ClientMock{GetWithPSKFunc: funcGetWithPsk}
// 	s3Clients := map[string]event.S3Client{bucket: s3cli}
// 	return s3cli, s3Clients
// }

// var contentLen = int64(284)

// // S3 GetWithPsk for a successful case
// var funcGetWithPskValid = func(key string, psk []byte) (io.ReadCloser, *int64, error) {
// 	return ioutil.NopCloser(strings.NewReader("exampleHeader\nexampleCsvLine")), &contentLen, nil
// }
