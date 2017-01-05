package main

import (
	"time"

	nsq "github.com/bitly/go-nsq"
)

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

// ReceiveMessage records the instrument about the ReceiveMessage function of the service.
func (mw InstrumentingMiddleware) ReceiveMessage(msg *nsq.Message) {
	defer func(begin time.Time) {
		lvs := []string{"method", "test", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.Service.ReceiveMessage(msg)
	return
}

func (service) ReceiveMessage(msg *nsq.Message) {
	//fmt.Println("Message received: " + string(msg.Body))
}
