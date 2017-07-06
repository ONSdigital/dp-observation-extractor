package mock

import (
	"github.com/ONSdigital/dp-observation-extractor/model"
	"github.com/ONSdigital/dp-observation-extractor/request"
)

var _ request.Handler = (*RequestHandler)(nil)

// NewRequestHandler returns a new mock request handler to capture request
func NewRequestHandler() *RequestHandler {

	requests := make([]model.Request, 0)

	return &RequestHandler{
		Requests: requests,
	}
}

// RequestHandler provides a mock implementation that captures requests to check.
type RequestHandler struct {
	Requests []model.Request
}

// Handle captures the given request and stores it for later assertions
func (handler *RequestHandler) Handle(request *model.Request) error {
	handler.Requests = append(handler.Requests, *request)
	return nil
}
