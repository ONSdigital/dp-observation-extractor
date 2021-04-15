package steps

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"github.com/ONSdigital/dp-observation-extractor/event"
	"github.com/ONSdigital/dp-observation-extractor/schema"
	"github.com/ONSdigital/dp-observation-extractor/service"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have a file in a bucket with name "([^"]*)" and content:$`, c.iHaveAFileInABucketWithNameAndContent)
	ctx.Step(`^I recieve a message containing the file name "([^"]*)"$`, c.iRecieveAMessageContainingTheFileName)
	ctx.Step(`^these messages are sent to the output message queue:$`, c.theseMessagesAreSentToTheOutputMessageQueue)
}

func (c *Component) iHaveAFileInABucketWithNameAndContent(filename string, content *godog.DocString) error {
	c.S3Client.GetWithPSKFunc = func(key string, psk []byte) (io.ReadCloser, *int64, error) {
		fmt.Println("=================================")
		fmt.Println(key)
		fmt.Println("=================================")
		length := int64(len(content.Content))
		return ioutil.NopCloser(strings.NewReader(content.Content)), &length, nil
	}
	return nil
}

func (c *Component) iRecieveAMessageContainingTheFileName(filename string) error {

	event := &event.DimensionsInserted{
		FileURL:    filename,
		InstanceID: "1",
	}

	signals := registerInterrupt()

	go func() {
		err := service.Run(context.Background(), c.cfg, c.serviceList, c.signals, c.errorChan, "", "", "")
		if err != nil {
			panic(err)
		}
	}()

	if err := c.sendToConsumer(event); err != nil {
		fmt.Println("=============================")
		return err
	}

	time.Sleep(3000 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func (c *Component) theseMessagesAreSentToTheOutputMessageQueue(table *godog.Table) error {
	return godog.ErrPending
}

func (c *Component) convertToEvents(table *godog.Table) ([]*event.DimensionsInserted, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&event.DimensionsInserted{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*event.DimensionsInserted), nil
}

func (c *Component) sendToConsumer(e *event.DimensionsInserted) error {
	bytes, err := schema.DimensionsInsertedEvent.Marshal(e)
	if err != nil {
		return err
	}

	c.KafkaConsumer.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	return nil

}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}
