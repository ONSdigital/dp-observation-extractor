package request

import (
	"github.com/ONSdigital/dp-observation-extractor/kafka"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/go-ns/log"
)

// MessageConsumer provides a generic interface for consuming []byte messages
type MessageConsumer interface {
	Incoming() chan kafka.Message
	Closer() chan bool
}

// Handler represents a handler for processing a single request.
type Handler interface {
	Handle(request *Request) error
}

// Consume convert them to request instances, and pass the request to the provided handler.
func Consume(messageConsumer MessageConsumer, handler Handler) {
	for message := range messageConsumer.Incoming() {

		request, err := Unmarshal(message)
		if err != nil {
			log.Error(err, log.Data{"schema": "Failed to unmarshal request"})
			continue
		}

		err = handler.Handle(request)
		if err != nil {
			log.Error(err, log.Data{"schema": "Failed to handle request"})
			continue
		}

		message.Commit()
	}
}

// Unmarshal converts a request instance to []byte.
func Unmarshal(message kafka.Message) (*Request, error) {
	var request Request
	err := schema.Request.Unmarshal(message.GetData(), &request)
	return &request, err
}
