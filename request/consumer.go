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
			log.Error(err, log.Data{"schema": "failed to unmarshal request"})
			continue
		}

		log.Debug("request received", log.Data{"request": request})

		err = handler.Handle(request)
		if err != nil {
			log.Error(err, log.Data{"schema": "failed to handle request"})
			continue
		}

		log.Debug("request processed - committing message", log.Data{"request": request})
		message.Commit()
		log.Debug("message committed", log.Data{"request": request})
	}
}

// Unmarshal converts a request instance to []byte.
func Unmarshal(message kafka.Message) (*Request, error) {
	var request Request
	err := schema.Request.Unmarshal(message.GetData(), &request)
	return &request, err
}
