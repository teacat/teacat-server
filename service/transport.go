package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// serviceHandlers returns the handlers that deal with the service.
func serviceHandlers(ctx context.Context, opts []httptransport.ServerOption, svc Service) []serviceHandler {

	uppercaseHandler := httptransport.NewServer(ctx, makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse, opts...)
	countHandler := httptransport.NewServer(ctx, makeCountEndpoint(svc), decodeCountRequest, encodeResponse, opts...)
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
			pattern: "/sd_health",
			handler: consulsdHandler,
		},
		{
			pattern: "/metrics",
			handler: stdprometheus.Handler(),
		},
	}
}

func eventListeners(svc Service) []eventListener {
	return []eventListener{
		{
			event:   "HelloWorld",
			body:    make(map[string]interface{}),
			meta:    make(map[string]string),
			handler: svc.CatchEvent,
		},
	}
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
