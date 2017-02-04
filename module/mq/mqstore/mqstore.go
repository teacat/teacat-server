package mqstore

import (
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	nsq "github.com/bitly/go-nsq"
)

type mqstore struct {
	*nsq.Producer
}

func NewProducer(producer string, httpProducer string, lookupds []string) *mqstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingMQ(httpProducer); err != nil {
		logrus.Fatalln(err)
	}

	config := nsq.NewConfig()
	prod, err := nsq.NewProducer(producer, config)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while creating the NSQ producer.")
	}

	return &mqstore{prod}
}

func pingMQ(addr string) error {
	for i := 0; i < 30; i++ {
		_, err := http.Get("http://" + addr)
		if err == nil {
			return nil
		}
		logrus.Infof("Cannot connect to NSQ producer, retry in 3 second.")
		time.Sleep(time.Second * 3)
	}

	return errors.New("Cannot connect to NSQ producer.")
}
