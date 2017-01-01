package main

import (
	"os"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/log"
)

// Uppercase 會紀錄 Uppercase 函式的相關資訊。
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

// Count 會紀錄 Count 函式的相關資訊。
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

//
//
//
//
//

// logmw is the middleware of the logger.
type LoggingMiddleware struct {
	Logger log.Logger
	Service
}

// loggingMiddleware recevie a logger and returns a new logging middleware.
func CreateLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return LoggingMiddleware{Logger: logger, Service: next}
	}
}

func createLogger(port *string) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", port).With("caller", log.DefaultCaller)

	return logger
}
