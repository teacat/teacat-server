package main

import (
	"flag"
	"net/http"

	"golang.org/x/net/context"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Svc struct {
	ConfigPath string
}

func main() {

	// Command line flags.
	var (
		listenPort = flag.String("listen", ":8080", "HTTP listen address")
	)

	flag.Parse()

	loadConfig("./")
	logger := createLogger(listenPort)
	db := createDatabase()
	model := createModel(db)
	msg := createMessage(logger)

	createService(
		logger,
		msg,
		model,
	)

	registerService(logger)

	logger.Log("msg", "HTTP", "addr", *listenPort)
	logger.Log("err", http.ListenAndServe(*listenPort, nil))
}

func createService(logger log.Logger, msg *nsq.Producer, model Model) {

	var svc Service
	svc = service{Message: msg, Model: model}
	svc = CreateLoggingMiddleware(logger)(svc)
	svc = CreateInstruMiddleware()(svc)

	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	setServiceHandlers(ctx, options, svc)
	setMessageHandlers(svc)

}

func setServiceHandlers(ctx context.Context, opts []httptransport.ServerOption, svc Service) {

	uppercaseHandler := httptransport.NewServer(ctx, makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse, opts...)
	countHandler := httptransport.NewServer(ctx, makeCountEndpoint(svc), decodeCountRequest, encodeResponse, opts...)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", stdprometheus.Handler())
}
