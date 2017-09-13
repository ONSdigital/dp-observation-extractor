package errors

import (
	"time"

	eventhandler "github.com/ONSdigital/dp-import-reporter/handler"
	eventSchema "github.com/ONSdigital/dp-import-reporter/schema"

	"github.com/ONSdigital/go-ns/log"
)

var _ Handler = (*KafkaHandler)(nil)

// Handler is a generic interface for error handling
type Handler interface {
	Handle(instanceID string, err error)
}

// KafkaHandler provides an error handler that writes to a kafka service
type KafkaHandler struct {
	messageProducer MessageProducer
}

// NewKafkaHandler returns a new KafkaHandler which will send messages to the producer
func NewKafkaHandler(messageProducer MessageProducer) *KafkaHandler {
	return &KafkaHandler{
		messageProducer: messageProducer,
	}
}

//MessageProducer dependency that writes messages
type MessageProducer interface {
	Output() chan []byte
	Closer() chan bool
}

// Handle logs the errors and sends it to the kafka error reporter
func (handler *KafkaHandler) Handle(instanceID string, err error) {
	data := log.Data{"INSTANCEID": instanceID, "ERROR": err.Error()}

	log.Info("Recieved error report", data)
	eventReport := eventhandler.EventReport{
		InstanceID: instanceID,
		EventType:  "error",
		EventMsg:   err.Error(),
	}
	errMsg, err := eventSchema.ReportedEventSchema.Marshal(eventReport)
	if err != nil {
		log.ErrorC("Failed to marshall error to event-reporter", err, data)
		return
	}
	handler.messageProducer.Output() <- errMsg

	time.Sleep(time.Duration(1000 * time.Millisecond))
}
