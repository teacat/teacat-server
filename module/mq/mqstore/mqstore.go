package mqstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
	config      *nsq.Config
	isConnected bool
	queue       []message
}

type message struct {
	topic string
	body  []byte
}

func NewProducer(url, producer, prodHTTP string, lookupds []string, m *mqutil.Engine, ready <-chan bool) *mqstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingMQ(prodHTTP); err != nil {
		logrus.Fatalln(err)
	}

	config := nsq.NewConfig()
	prod, err := nsq.NewProducer(producer, config)
	prod.SetLogger(nil, nsq.LogLevelError)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while creating the NSQ producer.")
	}

	ms := &mqstore{prod, config, true, []message{}}

	go ms.capture(url, prodHTTP, lookupds, m, ready)

	go ms.push(prodHTTP)

	return ms
}

// TODO: PING LOOKUPD

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

func (ms *mqstore) capture(url string, prodHTTP string, lookupds []string, m *mqutil.Engine, ready <-chan bool) {
	// Continue if the router was ready.
	<-ready

	for _, v := range m.Listeners {
		fmt.Println(v)
		c, err := nsq.NewConsumer(v.Topic, v.Channel, nsq.NewConfig())
		c.SetLogger(nil, nsq.LogLevelError)
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

func (mq *mqstore) push(prodHTTP string) {
	for {
		// Check the queue every second.
		<-time.After(time.Second * 1)

		//calculate send rateRATE

		// Ping the Event Store to see if it's back online or not.
		if !mq.isConnected {
			if err := pingMQ(prodHTTP); err == nil {
				mq.isConnected = true

				logrus.Infof("NSQ Producer is back online, there are %d unsent messages that will begin to send.", len(mq.queue))
			}
			continue
		}

		// Skip if there's nothing in the queue.
		if len(mq.queue) == 0 {
			continue
		}

		// A downward loop for the queue.
		for i := len(mq.queue) - 1; i >= 0; i-- {
			m := mq.queue[i]

			// Append the event in the stream.
			err := mq.Publish(m.topic, m.body)
			if err != nil {
				continue
			}

			// Remove the event from the queue since it has been sent.
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
		}
	}
}

func (mq *mqstore) send(topic string, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		//return err
	}
	if err := mq.Publish(topic, body); err != nil {
		switch t := err.(type) {
		case *net.OpError:
			// Mayne connect refuse
			if t.Op == "dial" {
				mq.isConnected = false
				mq.queue = append([]message{message{topic, body}}, mq.queue...)
				logrus.Warningf("The `%s` message will be sent when the NSQ Producer is back online. (Queue length: %d)", topic, len(mq.queue))
			}
		}
	}
}
