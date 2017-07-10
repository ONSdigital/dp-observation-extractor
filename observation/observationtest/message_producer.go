package observationtest

import "github.com/ONSdigital/dp-observation-extractor/observation"

var _ observation.MessageProducer = (*MessageProducer)(nil)

type MessageProducer struct {
	observation.MessageProducer
	OutputChannel chan []byte
}

func (messageProducer MessageProducer) Output() chan []byte {
	return messageProducer.OutputChannel
}
