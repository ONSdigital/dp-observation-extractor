package message

import (
	"github.com/ONSdigital/dp-observation-extractor/model"
	"github.com/ONSdigital/dp-observation-extractor/request"
	"github.com/ONSdigital/go-ns/avro"
	"github.com/ONSdigital/go-ns/log"
)

// Consumer provides a generic interface for consuming []byte messages
type Consumer interface {
	Messages() chan []byte
	Closer() chan bool
}

// ConsumeMessages convert them to request instances, and pass the request to the provided handler.
func ConsumeMessages(consumer Consumer, handler request.Handler) {
	for message := range consumer.Messages() {

		request, err := ToRequest(message)
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

// ToRequest converts the given []byte to a Request instance.
func ToRequest(message []byte) (*model.Request, error) {
	marshalSchema := &avro.Schema{
		Definition: RequestSchema,
	}

	var request model.Request
	err := marshalSchema.Unmarshal(message, &request)
	return &request, err
}
