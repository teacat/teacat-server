package main

import (
	"time"

	nsq "github.com/bitly/go-nsq"
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

// PublishMessage logs the informations about the PublishMessage function of the service.
func (mw LoggingMiddleware) ServiceDiscoveryCheck() {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "service_discovery_check",
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.ServiceDiscoveryCheck()
	return
}

// PublishMessage logs the informations about the PublishMessage function of the service.
func (mw LoggingMiddleware) PublishMessage(s string) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "publish_message",
			"input", s,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.PublishMessage(s)
	return
}

// ReceiveMessage logs the informations about the ReceiveMessage function of the service.
func (mw LoggingMiddleware) ReceiveMessage(msg *nsq.Message) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "receive_message",
			"input", msg.Body,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.Service.ReceiveMessage(msg)
	return
}
