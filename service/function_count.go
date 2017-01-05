package main

import (
	"time"
)

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

// Count records the instrument about the Count function of the service.
func (mw InstrumentingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(n))
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

// Count counts the length of the string.
func (svc service) Count(s string) int {
	return svc.Model.Count(s)
}
