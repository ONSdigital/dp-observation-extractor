package kafka

type Producer struct {
	// kafka producer
}

func (producer Producer) Output()   chan []byte {
	return nil
}

func (producer Producer) Closer()   chan bool {
	return nil
}