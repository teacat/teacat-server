package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

// loggingMiddleware recevie a logger and returns a new logging middleware.
func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next StringService) StringService {
		return logmw{logger, next}
	}
}

// logmw is the middleware of the logger.
type logmw struct {
	logger log.Logger
	StringService
}

// Uppercase 會紀錄 Uppercase 函式的相關資訊。
func (mw logmw) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.StringService.Uppercase(s)
	return
}

// Lowercase 會紀錄 Lowercase 函式的相關資訊。
func (mw logmw) Lowercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "lowercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.StringService.Lowercase(s)
	return
}

// Count 會紀錄 Count 函式的相關資訊。
func (mw logmw) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.StringService.Count(s)
	return
}
