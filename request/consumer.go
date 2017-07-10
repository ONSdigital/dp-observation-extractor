package request

import (
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/go-ns/log"
)

// MessageConsumer provides a generic interface for consuming []byte messages
type MessageConsumer interface {
	Messages() chan []byte
	Closer() chan bool
}

// Handler represents a handler for processing a single request.
type Handler interface {
	Handle(request *Request) error
}

// Consume convert them to request instances, and pass the request to the provided handler.
func Consume(messageConsumer MessageConsumer, handler Handler) {
	for message := range messageConsumer.Messages() {

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
	}
}

// Unmarshal converts a request instance to []byte.
func Unmarshal(input []byte) (*Request, error) {
	var request Request
	err := schema.Request.Unmarshal(input, &request)
	return &request, err
}
