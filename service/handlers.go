package main

import (
	"context"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// Create the service handlers with the serviceHandler struct:
//     serviceHandler{
//         pattern: "/uppercase",
//         handler: uppercaseHandler,
//     }

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

// Create the message handlers with the messageHandler struct:
//     messageHandler{
//         topic:   "new_user",
//         channel: "string",
//         handler: svc.NewUser,
//     }

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
