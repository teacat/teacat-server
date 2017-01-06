package main

import "time"

// Count logs the informations about the Count function of the service.
func (mw LoggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		mw.Logger.Log("method", "count", "input", s, "n", n, "took", time.Since(begin))
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

// Uppercase logs the informations about the Uppercase function of the service.
func (mw LoggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log("method", "uppercase", "input", s, "output", output, "err", err, "took", time.Since(begin))
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}
