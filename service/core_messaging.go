package main

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"

	nsq "github.com/bitly/go-nsq"
)

type messageHandlerFunc func(*nsq.Message)

type messageHandler struct {
	topic   string
	channel string
	handler messageHandlerFunc
}

func createMessage(l log.Logger) *nsq.Producer {
	prod, err := nsq.NewProducer(os.Getenv("KITSVC_NSQ_PRODUCER"), nsq.NewConfig())
	if err != nil {
		l.Log("module", "nsq", "msg", err)
	}

	return prod
}

func messageSubscribe(topic string, ch string, fn messageHandlerFunc) {

	q, err := nsq.NewConsumer(topic, ch, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	q.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		fn(msg)

		return nil
	}))

	if err := q.ConnectToNSQLookupds(strings.Split(os.Getenv("KITSVC_NSQ_LOOKUPS"), ",")); err != nil {
		panic(err)
	}
}

func setMessageSubscription(handlers []messageHandler) {
	for _, v := range handlers {
		messageSubscribe(v.topic, v.channel, v.handler)
	}
}
