package requesttest

import (
	"github.com/ONSdigital/dp-observation-extractor/request"
)

var _ request.MessageConsumer = (*MessageConsumer)(nil)

// NewMessageConsumer creates a consumer using the given channel.
func NewMessageConsumer(messages chan []byte) *MessageConsumer {
	return &MessageConsumer{
		messages: messages,
		closer:   nil,
	}
}

// MessageConsumer is a mock that provides the stored schema channel.
type MessageConsumer struct {
	messages chan []byte
	closer   chan bool
}

// Messages returns the stored schema channel.
func (consumer *MessageConsumer) Messages() chan []byte {
	return consumer.messages
}

// Closer returns the stored closer channel.
func (consumer *MessageConsumer) Closer() chan bool {
	return consumer.closer
}
