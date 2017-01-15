package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/TeaMeow/KitSvc/service/event"
	httptransport "github.com/go-kit/kit/transport/http"
)

// serviceHandlers returns the service handlers.
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

// eventListeners returns the event handlers.
func eventListeners(svc Service) []eventListener {
	return []eventListener{
		{
			event:   "uppercase",
			body:    &event.String{},
			meta:    make(H),
			handler: svc.CatchUppercase,
		},
	}
}

// decodeCountRequest decodes the request of the Count operation.
func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// decodeUppercaseRequest decodes the request of the Uppercase operation.
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
