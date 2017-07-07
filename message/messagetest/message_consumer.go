package messagetest

import "github.com/ONSdigital/dp-observation-extractor/message"

var _ message.Consumer = (*messageConsumer)(nil)

// NewMessageConsumer creates a consumer using the given channel.
func NewMessageConsumer(messages chan []byte) message.Consumer {
	return &messageConsumer{
		messages: messages,
		closer:   nil,
	}
}

type messageConsumer struct {
	messages chan []byte
	closer   chan bool
}

func (consumer *messageConsumer) Messages() chan []byte {
	return consumer.messages
}

func (consumer *messageConsumer) Closer() chan bool {
	return consumer.closer
}
