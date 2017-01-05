package main

import (
	"fmt"
	"time"
)

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

// Uppercase records the instrument about the Uppercase function of the service.
func (mw InstrumentingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

// Uppercase converts the string to uppercase.
func (svc service) Uppercase(s string) (string, error) {

	res, err := svc.Model.ToUpper(s)
	if err != nil {
		return "", err
	}

	return res, nil
}
