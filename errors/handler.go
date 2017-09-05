package errors

import (
	"time"

	eventhandler "github.com/ONSdigital/dp-import-reporter/handler"
	eventSchema "github.com/ONSdigital/dp-import-reporter/schema"

	"github.com/ONSdigital/go-ns/log"
)

type Handler interface {
	Handle(instanceID string, err error, data log.Data)
}
type KafkaHandler struct {
	messageProducer MessageProducer
}

func NewKafkaHandler(messageProducer MessageProducer) *KafkaHandler {
	return &KafkaHandler{
		messageProducer: messageProducer,
	}
}

type MessageProducer interface {
	Output() chan []byte
	Closer() chan bool
}

func (handler *KafkaHandler) Handle(instanceID string, err error, data log.Data) {
	eventReport := eventhandler.EventReport{
		InstanceID: instanceID,
		EventType:  "error",
		EventMsg:   err.Error(),
	}
	errMsg, err := eventSchema.ReportedEventSchema.Marshal(eventReport)
	if err != nil {
		log.Error(err, nil)
		return
	}
	handler.messageProducer.Output() <- errMsg

	time.Sleep(time.Duration(1000 * time.Millisecond))
}
