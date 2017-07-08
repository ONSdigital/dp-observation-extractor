package observation

// MessageWriter writes observations as messages
type MessageWriter struct {
	messageProducer MessageProducer
}

// MessageProducer represents the dependency that writes messages
type MessageProducer interface{
	Output()   chan []byte
	Closer()   chan bool
}

// NewMessageWriter returns a new message writer that writes messages to the given message producer.
func NewMessageWriter(messageProducer MessageProducer) *MessageWriter {
	return &MessageWriter{
		messageProducer: messageProducer,
	}
}

// WriteAll observations from the given observation reader as messages.
func (messageWriter MessageWriter) WriteAll(reader Reader) error {
	return nil
}

