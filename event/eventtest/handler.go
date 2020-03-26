package eventtest

import (
	"context"

	"github.com/ONSdigital/dp-observation-extractor/event"
)

var _ event.Handler = (*EventHandler)(nil)

// NewEventHandler returns a new mock event handler to capture event, with an optional error to return on Handle
func NewEventHandler(err error) *EventHandler {
	return &EventHandler{
		Events:   make([]event.DimensionsInserted, 0),
		ChHandle: make(chan *event.DimensionsInserted),
		Error:    err,
	}
}

// EventHandler provides a mock implementation that captures events to check.
type EventHandler struct {
	Events   []event.DimensionsInserted
	Error    error
	ChHandle chan *event.DimensionsInserted
}

// Handle captures the given event and stores it for later assertions
func (handler *EventHandler) Handle(ctx context.Context, event *event.DimensionsInserted) error {
	handler.Events = append(handler.Events, *event)
	handler.ChHandle <- event
	return handler.Error
}
