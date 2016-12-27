package service

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

func SetHandlers(svc Service) {
	ctx := context.Background()

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)
	lowercaseHandler := httptransport.NewServer(
		ctx,
		makeLowercaseEndpoint(svc),
		decodeLowercaseRequest,
		encodeResponse,
	)
	countHandler := httptransport.NewServer(
		ctx,
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	http.Handle("/lowercase", lowercaseHandler)
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", stdprometheus.Handler())
}
