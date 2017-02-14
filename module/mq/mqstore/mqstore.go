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

var (
	// AllConnected returns true when the message queue were all connected.
	AllConnected = false
	// SentTotal returns the total of the sent message.
	SentTotal = 0
	// RecvTotal returns the total of the received message.
	RecvTotal = 0
	// QueueTotal returns the total of the message are still in the queue.
	QueueTotal = 0
)

type mqstore struct {
	*nsq.Producer
	config      *nsq.Config
	isConnected bool
	queue       []message
}

// message represents a message.
type message struct {
	topic string
	body  []byte
}

// NewProducer creates a new NSQ producer, and start to capturing the incoming messages.
func NewProducer(url, producer, prodHTTP string, lookupds []string, m *mqutil.Engine, deployed <-chan bool) *mqstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingMQ(prodHTTP); err != nil {
		logrus.Fatalln(err)
	}

	config := nsq.NewConfig()
	prod, err := nsq.NewProducer(producer, config)
	//prod.SetLogger(nil, nsq.LogLevelError)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while creating the NSQ producer.")
	}

	ms := &mqstore{
		Producer:    prod,
		config:      config,
		isConnected: true,
	}

	// Capturing the messages when the router was ready in the goroutine.
	go ms.capture(url, prodHTTP, lookupds, m, deployed)
	// Pushing the messages which are in the local queue to the remote message queue.
	go ms.push(prodHTTP)

	return ms
}

// TODO: PING LOOKUPD

// pingMQ pings the NSQ with backoff to ensure
// a connection can be established before we proceed with the
// message queue setup and migration.
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

// createTopic creates the new topic in the remote message queue,
// so we can subscribe to it.
func createTopic(httpProducer, topic string) {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/topic/create?topic=%s", httpProducer, topic))
	cmd.Start()
	cmd.Wait()
}

type logger struct {
}

func (l *logger) Output(calldepth int, s string) error {
	logger := logrus.StandardLogger()
	typ := s[0:3]
	switch typ {
	case "DBG":
		logger.Debug(s[9:len(s)])
	case "INF":
		logger.Info(s[9:len(s)])
	case "WRN":
		logger.Warn(s[9:len(s)])
	case "ERR":
		logger.Error(s[9:len(s)])
	}

	return nil
}

// capture the incoming events.
func (mq *mqstore) capture(url string, prodHTTP string, lookupds []string, m *mqutil.Engine, deployed <-chan bool) {
	// Continue if the router was ready.
	<-deployed
	// Each of the topic listener.
	for _, v := range m.Listeners {
		c, err := nsq.NewConsumer(v.Topic, v.Channel, nsq.NewConfig())
		l := &logger{}
		c.SetLogger(l, nsq.LogLevelDebug)
		if err != nil {
			logrus.Errorln(err)
			logrus.Fatalf("Cannot create the NSQ `%s` consumer. (channel: %s)", v.Topic, v.Channel)
		}

		// Create the topic to make sure it does exist before we subscribe to it.
		createTopic(prodHTTP, v.Topic)
		// Add the topic handler.
		c.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
			RecvTotal++
			// Send the received message to the self router,
			// so we can process it with Gin.
			sendToRouter(v.Method, url+v.Path, msg.Body)
			return nil
		}))

		// Connect to the NSQLookupds instead of a single NSQ node.
		if err := c.ConnectToNSQLookupds(lookupds); err != nil {
			logrus.Errorln(err)
			logrus.Fatalln("Cannot connect to the NSQ lookupds.")
		}
	}
}

// push the message which are in the queue to the remote message queue.
func (mq *mqstore) push(prodHTTP string) {
	for {
		<-time.After(time.Millisecond * 10)

		// Ping the NSQ Producer to see if it's back online or not.
		if !mq.isConnected {
			if err := pingMQ(prodHTTP); err == nil {
				mq.isConnected = true
				AllConnected = true

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
			// Wait a little bit for another event.
			<-time.After(time.Millisecond * 2)
			// Append the message in the topic.
			err := mq.Publish(m.topic, m.body)
			if err != nil {
				continue
			}
			QueueTotal--
			// Remove the message from the queue since it has been sent.
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
		}
	}
}

// send the event to the specified stream.
func (mq *mqstore) send(topic string, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		//return err
	}
	// Counter.
	SentTotal++

	// Send the message to the remote message queue.
	if err := mq.Publish(topic, body); err != nil {
		switch t := err.(type) {
		case *net.OpError:
			// Push the message to the local queue if connecting refused.
			if t.Op == "dial" {
				// Mark the connection as lost.
				mq.isConnected = false
				AllConnected = false
				// Push the message to local queue.
				mq.queue = append([]message{message{topic, body}}, mq.queue...)
				// Counter.
				QueueTotal++
				logrus.Warningf("The `%s` message will be sent when the NSQ Producer is back online. (queue length: %d)", topic, len(mq.queue))
			}
		}
	}
}
