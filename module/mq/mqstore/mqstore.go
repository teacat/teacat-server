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
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/module/mq"
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

// lg is a logger for NSQ.
type lg struct {
}

// NewProducer creates a new NSQ producer, and start to capturing the incoming messages.
func NewProducer(url, producer, prodHTTP string, lookupds []string, m *mqutil.Engine, deployed <-chan bool) *mqstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingMQ(prodHTTP); err != nil {
		logger.Fatal(err)
	}
	config := nsq.NewConfig()
	// Create the producer.
	prod, err := nsq.NewProducer(producer, config)
	// Set the logger for the producer.
	prod.SetLogger(&lg{}, nsq.LogLevelDebug)
	if err != nil {
		logger.FatalFields("Error occurred while creating the NSQ producer.", logrus.Fields{
			"err":      err,
			"url":      url,
			"producer": producer,
			"http":     prodHTTP,
			"lookupds": lookupds,
		})
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

// Output the error to the global logger.
func (l *lg) Output(calldepth int, s string) error {
	typ := s[0:3]
	switch typ {
	case "DBG":
		logger.Debug(s[9:len(s)])
	case "INF":
		// Use debug so it won't be shown in the terminal.
		logger.Debug(s[9:len(s)])
	case "WRN":
		logger.Warning(s[9:len(s)])
	case "ERR":
		logger.Error(s[9:len(s)])
	}
	return nil
}

// pingMQ pings the NSQ with backoff to ensure
// a connection can be established before we proceed with the
// message queue setup and migration.
func pingMQ(addr string) error {
	for i := 0; i < 30; i++ {
		_, err := http.Get("http://" + addr)
		if err == nil {
			return nil
		}
		logger.Info("Cannot connect to NSQ producer, retry in 1 second.")
		time.Sleep(time.Second * 1)
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
		logger.ErrorFields("Error occurred while sending the message to self router.", logrus.Fields{
			"err":    err,
			"method": method,
			"url":    url,
			"body":   json,
		})
	}
	if resp.StatusCode != 200 {
		logger.InfoFields("The message has been recevied by the router, but the status code wasn't 200.", logrus.Fields{
			"status": resp.StatusCode,
			"method": method,
			"url":    url,
			"body":   json,
		})
	}
}

// createTopic creates the new topic in the remote message queue,
// so we can subscribe to it.
func createTopic(httpProducer, topic string) {
	cmd := exec.Command("curl", "-X", "POST", fmt.Sprintf("http://%s/topic/create?topic=%s", httpProducer, topic))
	cmd.Start()
	cmd.Wait()
}

// capture the incoming events.
func (mq *mqstore) capture(url string, prodHTTP string, lookupds []string, m *mqutil.Engine, deployed <-chan bool) {
	// Continue if the router was ready.
	<-deployed
	// Each of the topic listener.
	for _, v := range m.Listeners {
		// Create the consumer.
		c, err := nsq.NewConsumer(v.Topic, v.Channel, nsq.NewConfig())
		// Set the logger for the consumer.
		c.SetLogger(&lg{}, nsq.LogLevelDebug)
		if err != nil {
			logger.FatalFields("Cannot create the NSQ consumer.", logrus.Fields{
				"channel": v.Channel,
				"topic":   v.Topic,
				"path":    v.Path,
				"method":  v.Method,
			})
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
			logger.FatalFields("Cannot connect to the NSQ lookupds.", logrus.Fields{
				"err":      err,
				"lookupds": lookupds,
				"channel":  v.Channel,
				"topic":    v.Topic,
			})
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

				logger.InfoFields("NSQ producer is back online, the unsent messages that will begin to send.", logrus.Fields{
					"unsent": len(mq.queue),
				})
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
			err := mq.Producer.Publish(m.topic, m.body)
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
		logger.ErrorFields("Error occurred while converting the data to json.", logrus.Fields{
			"err":   err,
			"topic": topic,
			"data":  data,
		})
	}
	// Counter.
	SentTotal++

	// Send the message to the remote message queue.
	if err := mq.Producer.Publish(topic, body); err != nil {
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
				logger.WarningFields("Message will be sent when the NSQ producer is back online.", logrus.Fields{
					"topic":  topic,
					"unsent": len(mq.queue),
				})
			}
		}
	}
}

// Publish the message.
func (mq *mqstore) Publish(m mq.M) {
	go mq.send(m.Topic, m.Data)
}
