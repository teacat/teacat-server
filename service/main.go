package main

import (
	"flag"
	"net/http"

	"golang.org/x/net/context"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	// Command line flags.
	var (
		listenPort = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	// Load the configurations.
	loadConfig("./")

	// Create the logger with the specified listen port.
	logger := createLogger(listenPort)
	// Create the database connection.
	db := createDatabase()
	// Create the model with the database connection.
	model := createModel(db)
	// Create the messaging service with the logger.
	msg := createMessage(logger)

	// Create the main service with what it needs.
	createService(logger, msg, model)
	// Register the service to the service registry.
	registerService(logger)

	// Log the ports.
	logger.Log("msg", "HTTP", "addr", *listenPort)
	// Start the service and listening to the requests.
	logger.Log("err", http.ListenAndServe(*listenPort, nil))
}

// createService creates the main service by setting the handlers and preparing the middlewares.
func createService(logger log.Logger, msg *nsq.Producer, model Model) {

	var svc Service
	svc = service{Message: msg, Model: model}
	svc = createLoggingMiddleware(logger)(svc)
	svc = createInstruMiddleware()(svc)

	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	setServiceSubscription(serviceHandlers(ctx, options, svc))
	setMessageSubscription(messageHandlers(svc))
}

type serviceHandler struct {
	pattern string
	handler http.Handler
}

func setServiceSubscription(handlers []serviceHandler) {
	for _, v := range handlers {
		http.Handle(v.pattern, v.handler)
	}
}
