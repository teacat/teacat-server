package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

type publishMessageRequest struct {
	S string `json:"s"`
}

type publishMessageResponse struct {
	V string `json:"v"`
}

func makePublishMessageEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(publishMessageRequest)
		svc.PublishMessage(req.S)
		return publishMessageResponse{req.S}, nil
	}
}

func decodePublishMessageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request publishMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// PublishMessage logs the informations about the PublishMessage function of the service.
func (mw LoggingMiddleware) PublishMessage(s string) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "publish_message",
			"input", s,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.PublishMessage(s)
	return
}

// PublishMessage records the instrument about the PublishMessage function of the service.
func (mw InstrumentingMiddleware) PublishMessage(s string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "publish_message", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.Service.PublishMessage(s)
	return
}

func (svc service) PublishMessage(s string) {
	svc.Message.Publish("hello_world", []byte(s))
}
