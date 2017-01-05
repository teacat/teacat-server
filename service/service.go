package main

import (
	"context"
	"errors"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	nsq "github.com/bitly/go-nsq"
)

// Error codes returned by failures
var (

	// ErrEmpty will returned if the string is empty.
	ErrEmpty = ErrInfo{
		Text:   errors.New("The string is empty."),
		Status: http.StatusBadRequest,
		Code:   "str_empty",
	}
)

// Service represents the operations of the service can do.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	PublishMessage(string)
	ReceiveMessage(*nsq.Message)
}

// serviceHandlers returns the handlers that deal with the service.
func serviceHandlers(ctx context.Context, opts []httptransport.ServerOption, svc Service) []serviceHandler {

	uppercaseHandler := httptransport.NewServer(ctx, makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse, opts...)
	countHandler := httptransport.NewServer(ctx, makeCountEndpoint(svc), decodeCountRequest, encodeResponse, opts...)
	publishHandler := httptransport.NewServer(ctx, makePublishMessageEndpoint(svc), decodePublishMessageRequest, encodeResponse, opts...)
	consulsdHandler := httptransport.NewServer(ctx, makeServiceDiscoveryEndpoint(svc), decodeServiceDiscoveryRequest, encodeResponse, opts...)

	return []serviceHandler{
		{
			pattern: "/uppercase",
			handler: uppercaseHandler,
		},
		{
			pattern: "/count",
			handler: countHandler,
		},
		{
			pattern: "/publish",
			handler: publishHandler,
		},
		{
			pattern: "/sd_health",
			handler: consulsdHandler,
		},
		{
			pattern: "/metrics",
			handler: stdprometheus.Handler(),
		},
	}
}

// messageHandlers returns the handlers that deal with the messages.
func messageHandlers(svc Service) []messageHandler {

	return []messageHandler{
		{
			topic:   "hello_world",
			channel: "string",
			handler: svc.ReceiveMessage,
		},
	}
}
