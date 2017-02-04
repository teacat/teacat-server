package mqstore

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/shared/mqutil"
	nsq "github.com/bitly/go-nsq"
	"github.com/parnurzeal/gorequest"
)

type mqstore struct {
	*nsq.Producer
}

func NewProducer(url, producer, httpProducer string, lookupds []string, m *mqutil.Engine) *mqstore {
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

// sendToRouter sends the received event data to self router.
func sendToRouter(method string, url string, json []byte) {
	// Send the request via the HTTP client.
	resp, _, err := gorequest.
		New().
		CustomMethod(method, url).
		Send(string(json)).
		End()
	if err != nil {
		logrus.Errorln(err)
		// not fatal, TOO PANIC!
		logrus.Fatalln("Error occurred while sending the event to self router.")
	}
	if resp.StatusCode != 200 {
		logrus.Infoln("The event has been recevied by the router, but the status code wasn't 200.")
	}
}

func createTopic(topic, httpProducer string) {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/topic/create?topic=%s", httpProducer, topic))
	cmd.Start()
	cmd.Wait()
}

func (ms *mqstore) capture(url string, prodHTTP string, lookupds []string, m *mqutil.Engine) {

	for _, v := range m.Listeners {
		c, err := nsq.NewConsumer(v.Topic, v.Channel, nsq.NewConfig())
		if err != nil {
			logrus.Errorln(err)
			logrus.Fatalf("Cannot create the NSQ `%s` consumer. (channel: %s)", v.Topic, v.Channel)
		}
		c.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
			//
			createTopic(prodHTTP, v.Path)
			//
			sendToRouter(v.Method, url+v.Path, msg.Body)

			return nil
		}))

		if err := c.ConnectToNSQLookupds(lookupds); err != nil {
			logrus.Errorln(err)
			logrus.Fatalln("Cannot connect to the NSQ lookupds.")
		}
	}

}
