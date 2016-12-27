package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/TeaMeow/KitSvc/config"
	"github.com/TeaMeow/KitSvc/instrumenting"
	"github.com/TeaMeow/KitSvc/logging"
	"github.com/TeaMeow/KitSvc/service"

	"golang.org/x/net/context"

	"github.com/TeaMeow/KitSvc/sd"
	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	conf := config.Load()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	ctx := context.Background()

	var svc service.Service

	svc = service.Concrete{}
	svc = logging.CreateMiddleware(logger)(svc)
	svc = instrumenting.CreateMiddleware(conf)(svc)

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

	sd.Register(&conf, logger)

	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
