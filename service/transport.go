package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	httptransport "github.com/go-kit/kit/transport/http"
)

// serviceHandlers returns the handlers that deal with the service.
func serviceHandlers(ctx context.Context, opts []httptransport.ServerOption, svc Service) []serviceHandler {

	uppercaseHandler := httptransport.NewServer(ctx, makePostUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse, opts...)
	countHandler := httptransport.NewServer(ctx, makePostCountEndpoint(svc), decodeCountRequest, encodeResponse, opts...)

	return []serviceHandler{
		{
			method:  "POST",
			pattern: "/uppercase",
			handler: uppercaseHandler,
		},
		{
			method:  "POST",
			pattern: "/count",
			handler: countHandler,
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
