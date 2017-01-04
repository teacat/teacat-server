package main

import (
	"fmt"
	"time"

	nsq "github.com/bitly/go-nsq"
)

// Instrumenting function measures the performance of the operations of the service.
//
// Let's say that you have a `Uppercase` operation,
// then you would have to create a instrumenting function for the `Uppercase` operation.
//
// Create the instrumenting functions with the following format:
//     func (mw InstrumentingMiddleware)...

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

// Test records the instrument about the Test function of the service.
func (mw InstrumentingMiddleware) Test(msg *nsq.Message) {
	defer func(begin time.Time) {
		lvs := []string{"method", "test", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.Service.Test(msg)
	return
}
