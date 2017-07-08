package request

import (
	"github.com/ONSdigital/go-ns/avro"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/dp-observation-extractor/message"
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
			log.Error(err, log.Data{"message": "Failed to unmarshal request"})
			continue
		}

		err = handler.Handle(request)
		if err != nil {
			log.Error(err, log.Data{"message": "Failed to handle request"})
			continue
		}
	}
}

// Unmarshal converts the given []byte to a Request instance.
func Unmarshal(input []byte) (*Request, error) {
	marshalSchema := &avro.Schema{
		Definition: message.RequestSchema,
	}

	var request Request
	err := marshalSchema.Unmarshal(input, &request)
	return &request, err
}
