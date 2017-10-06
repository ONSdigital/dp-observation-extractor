package client

import (
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-reporter-client/model"
	"github.com/ONSdigital/dp-reporter-client/schema"
	"github.com/ONSdigital/go-ns/log"
)

const (
	instanceEmpty    = "cannot ReportError, instanceID is a required field but was empty"
	contextEmpty     = "cannot ReportError, errContext is a required field but was empty"
	sendingEvent     = "sending reportEvent for application error"
	failedToMarshal  = "failed to marshal reportEvent to avro"
	eventMessageFMT  = "%s: %s"
	serviceNameEmpty = "cannot create new reporter client as serviceName is empty"
	kafkaProducerNil = "cannot create new reporter client as kafkaProducer is nil"
	eventTypeErr     = "error"
	reportEventKey   = "reportEvent"
)

// KafkaProducer interface of the go-ns kafka.Producer
type KafkaProducer interface {
	Output() chan []byte
}

type marshalFunc func(s interface{}) ([]byte, error)

// ReporterClient a client for sending error reports to the import-reporter
type ReporterClient struct {
	kafkaProducer KafkaProducer
	marshal       marshalFunc
	serviceName   string
}

// NewReporterClientMock create a new ReporterClient to send error reports to the import-reporter
func NewReporterClient(kafkaProducer KafkaProducer, serviceName string) (*ReporterClient, error) {
	if kafkaProducer == nil {
		return nil, errors.New(kafkaProducerNil)
	}
	if len(serviceName) == 0 {
		return nil, errors.New(serviceNameEmpty)
	}
	cli := &ReporterClient{
		serviceName:   serviceName,
		kafkaProducer: kafkaProducer,
		marshal:       schema.ReportEventSchema.Marshal,
	}
	return cli, nil
}

// ReportError send an error report to the import-reporter
func (c ReporterClient) ReportError(instanceID string, errContext string, err error, data log.Data) error {
	if len(instanceID) == 0 {
		log.Info(instanceEmpty, nil)
		return errors.New(instanceEmpty)
	}
	if len(errContext) == 0 {
		log.Info(contextEmpty, nil)
		return errors.New(contextEmpty)
	}

	reportEvent := &model.ReportEvent{
		InstanceID:  instanceID,
		EventMsg:    fmt.Sprintf(eventMessageFMT, errContext, err.Error()),
		ServiceName: c.serviceName,
		EventType:   eventTypeErr,
	}

	reportEventData := log.Data{reportEventKey: reportEvent}
	log.Info(sendingEvent, reportEventData)

	avroBytes, err := c.marshal(reportEvent)
	if err != nil {
		log.ErrorC(failedToMarshal, err, reportEventData)
		return err
	}

	c.kafkaProducer.Output() <- avroBytes
	return nil
}
