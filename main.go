package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/TeaMeow/KitSvc/config"
	"github.com/TeaMeow/KitSvc/discovery"
	"github.com/TeaMeow/KitSvc/instrumenting"
	"github.com/TeaMeow/KitSvc/logging"
	"github.com/TeaMeow/KitSvc/messaging"
	"github.com/TeaMeow/KitSvc/service"

	"github.com/go-kit/kit/log"
)

func main() {

	// Command line flags.
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	// Load the configurations.
	conf := config.Load()

	// The logger.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)

	// Register to the message system.
	msg := messaging.Create(conf, logger)

	// Service and the middlewares.
	var svc service.Service
	svc = service.Concrete{Message: msg.Producer}
	svc = logging.CreateMiddleware(logger)(svc)
	svc = instrumenting.CreateMiddleware(conf)(svc)

	//
	msg.SetHandlers(svc)

	// Set the handlers.
	service.SetHandlers(svc)

	// Register the service to the service discovery registry.
	discovery.Register(&conf, logger)

	// Log and start the HTTP transmission.
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
