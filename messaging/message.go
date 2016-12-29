package messaging

import (
	"github.com/go-kit/kit/log"

	"github.com/TeaMeow/KitSvc/config"
	nsq "github.com/bitly/go-nsq"
)

type Concrete struct {
	Producer *nsq.Producer
	Config   config.Context
	Logger   log.Logger
}

func Create(c config.Context, logger log.Logger) Concrete {

	prod, err := nsq.NewProducer(c.NSQ.Producer, nsq.NewConfig())
	if err != nil {
		logger.Log("module", "nsq", "msg", err)
	}

	return Concrete{Producer: prod, Config: c}
}

type HandlerFunc func(*nsq.Message)

func (c Concrete) Handle(topic string, ch string, fn HandlerFunc) {

	q, err := nsq.NewConsumer(topic, ch, nsq.NewConfig())
	if err != nil {
		c.Logger.Log("module", "nsq", "msg", err)
	}

	q.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		fn(msg)

		return nil
	}))

	if err := q.ConnectToNSQLookupds(c.Config.NSQ.Lookups); err != nil {
		c.Logger.Log("module", "nsq", "msg", err)
	}
}
