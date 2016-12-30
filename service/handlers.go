package service

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

func SetHandlers(svc Service) {
	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
		options...,
	)
	lowercaseHandler := httptransport.NewServer(
		ctx,
		makeLowercaseEndpoint(svc),
		decodeLowercaseRequest,
		encodeResponse,
		options...,
	)
	countHandler := httptransport.NewServer(
		ctx,
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
		options...,
	)

	http.Handle("/lowercase", lowercaseHandler)
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", stdprometheus.Handler())
}
