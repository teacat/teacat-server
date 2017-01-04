package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

// LoggingMiddleware represents a middleware of the logger.
type LoggingMiddleware struct {
	Logger log.Logger
	Service
}

// createLoggingMiddleware creates the logging middleware.
func createLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return LoggingMiddleware{Logger: logger, Service: next}
	}
}

// createLogger creates the logger with the specified port which tracks the function callers.
func createLogger(port *string) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", port).With("caller", log.DefaultCaller)

	return logger
}
