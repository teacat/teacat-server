package main

import (
	"os"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/log"
)

// Logging function logs the input, output and the caller of the operations of the service.
//
// Let's say that you have a `Uppercase` operation,
// then you would have to create a logging function for the `Uppercase` operation.
//
// Create the logging functions with the following format:
//     func (mw LoggingMiddleware)...

// Uppercase logs the informations about the Uppercase function of the service.
func (mw LoggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

// Count logs the informations about the Count function of the service.
func (mw LoggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

// Test logs the informations about the Test function of the service.
func (mw LoggingMiddleware) Test(msg *nsq.Message) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "test",
			"input", msg.Body,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.Test(msg)
	return
}

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
