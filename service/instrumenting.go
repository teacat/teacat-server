package main

import (
	"fmt"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
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

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

// InstrumentingMiddleware represents a middleware of the instrumenting.
type InstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	Service
}

// createInstruMiddleware creates the instrumenting middleware.
func createInstruMiddleware() ServiceMiddleware {

	fieldKeys := []string{"method", "error"}

	// Number of requests received.
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	// Total duration of requests in microseconds.
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// The result of each count method.
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: viper.GetString("prometheus.namespace"),
		Subsystem: viper.GetString("prometheus.subsystem"),
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next Service) Service {
		return InstrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}
