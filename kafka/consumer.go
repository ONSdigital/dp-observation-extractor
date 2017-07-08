package kafka

type Consumer struct {
	//kafka.Consumer
}

func (consumer Consumer) Messages() chan []byte {
	return nil
}
func (consumer Consumer) Closer() chan bool {
	return nil
}
