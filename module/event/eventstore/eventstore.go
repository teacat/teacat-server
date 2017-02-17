package eventstore

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	evt "github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/logger"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/TeaMeow/KitSvc/version"
	"github.com/jetbasrawi/go.geteventstore"
	"github.com/parnurzeal/gorequest"
)

var (
	// AllConnected returns true when the event stores were all connected.
	AllConnected = false
	// SentTotal returns the total of the sent event.
	SentTotal = 0
	// RecvTotal returns the total of the received event.
	RecvTotal = 0
	// QueueTotal returns the total of the event are still in the queue.
	QueueTotal = 0
)

type eventstore struct {
	*goes.Client
	isConnected bool
	queue       []event
}

// event represents an event data.
type event struct {
	stream string
	data   interface{}
	meta   interface{}
}

// NewClient creates a new event store client.
func NewClient(url string, esURL string, username string, password string, e *eventutil.Engine, replayed chan<- bool, deployed <-chan bool) *eventstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingStore(esURL); err != nil {
		logger.Fatal(err)
	}
	// Create the client to Event Store.
	client, err := goes.NewClient(nil, esURL)
	if err != nil {
		logger.FatalFields("Cannot create the client for Event Store.", logrus.Fields{
			"err":    err,
			"remote": esURL,
		})
	}
	client.SetBasicAuth(username, password)

	es := &eventstore{
		Client:      client,
		isConnected: true,
	}
	AllConnected = true

	// Capturing the events when the router was ready in the goroutine.
	go es.capture(url, e, replayed, deployed)
	// Pushing the events which are in the local queue to Event Store.
	go es.push(esURL)

	return es
}

// pingStore pings the eventstore with backoff to ensure
// a connection can be established before we proceed with the
// eventstore setup and migration.
func pingStore(url string) error {
	for i := 0; i < 30; i++ {
		_, err := http.Get(url)
		if err == nil {
			return nil
		}
		logger.Info("Cannot connect to Event Store, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to Event Store.")
}

// sendToRouter sends the received event data to the self router.
func sendToRouter(method string, url string, json []byte) {
	// Send the request via the HTTP client.
	resp, _, err := gorequest.
		New().
		CustomMethod(method, url).
		Send(string(json)).
		End()
	if err != nil {
		logger.ErrorFields("Error occurred while sending the event to self router.", logrus.Fields{
			"err":    err,
			"method": method,
			"url":    url,
			"body":   json,
		})
	}
	if resp.StatusCode != 200 {
		logger.InfoFields("The event has been recevied by the router, but the status code wasn't 200.", logrus.Fields{
			"status": resp.StatusCode,
			"method": method,
			"url":    url,
			"body":   json,
		})
	}
}

// isEmptyEvent returns true when the event body is empty.
func isEmptyEvent(json []byte) bool {
	return string(json) == "{}"
}

// capture the incoming events.
func (es *eventstore) capture(localURL string, e *eventutil.Engine, replayed chan<- bool, deployed <-chan bool) {
	<-deployed

	playedStreams := 0
	totalStreams := len(e.Listeners)

	// Each of the listener.
	for _, l := range e.Listeners {
		r := es.NewStreamReader(l.Stream)

		go func(l eventutil.Listener) {
			replayedStream := false
			startTime := time.Now()
			totalEvents := 0

			// Read the next event.
			for r.Next() {
				// Error occurred.
				if r.Err() != nil {
					switch r.Err().(type) {

					// Continue if there's no more event.
					case *goes.ErrNoMoreEvents:
						// We've replayed all the events in the stream since there're no more events can read,
						// so we mark the stream as `replayed`, we'll continue to the next step once all the streams were replayed.
						if !replayedStream {
							// Mark the stream has been replayed.
							playedStreams++
							replayedStream = true
							// Show the statistic of the replayed stream.
							logger.InfoFields("The event stream has been replayed successfully.", logrus.Fields{
								"stream": l.Stream,
								"wasted": time.Since(startTime).Seconds(),
								"amount": totalEvents,
							})
							// Set the `played` as true once we have replayed all the event streams.
							if playedStreams >= totalStreams {
								close(replayed)
							}
						}
						// When there are no more events in the stream, set LongPoll.
						// The server will wait for 15 seconds in this case or until
						// events become available on the stream.
						r.LongPoll(15)

					// Create an empty event if the stream hasn't been created.
					case *goes.ErrNotFound:
						writer := es.NewStreamWriter(l.Stream)

						err := writer.Append(nil, goes.NewEvent("", "", map[string]string{}, map[string]string{}))
						if err != nil {
							logger.FatalFields("Error occurred while creating an empty stream.", logrus.Fields{
								"stream": l.Stream,
								"err":    err,
							})
						}
						continue

						// Sleep for few seconds and try again if the Event Store was not connected.
					case *url.Error, *goes.ErrTemporarilyUnavailable:
						es.isConnected = false
						AllConnected = false
						<-time.After(time.Duration(3) * time.Second)

						// Other errors.
					default:
						logger.ErrorFields("Error occurred while reading the incoming event.", logrus.Fields{
							"stream": l.Stream,
							"err":    r.Err(),
						})
						continue
					}

					// We received the event, and we're going to make a http request,
					// abd send it to out own Gin router.
				} else {
					// Get the event body.
					json, err := r.EventResponse().Event.Data.(*json.RawMessage).MarshalJSON()
					if err != nil {
						logger.ErrorFields("Cannot parse the event data, the event has been skipped.", logrus.Fields{
							"stream": l.Stream,
							"err":    err,
						})
						continue
					}
					// Skip the empty json event,
					// because that might be the one which we used to create the new stream.
					if isEmptyEvent(json) {
						continue
					}
					// Counters.
					totalEvents++
					RecvTotal++
					// Send the received data to the router,
					// so we can process the event in the router.
					go sendToRouter(l.Method, localURL+l.Path, json)

					continue
				}
			}
		}(l)
	}
}

// push the events which are in the queue to the event store.
func (es *eventstore) push(esURL string) {
	for {
		<-time.After(time.Millisecond * 10)

		// Ping the Event Store to see if it's back online or not.
		if !es.isConnected {
			if err := pingStore(esURL); err == nil {
				es.isConnected = true
				AllConnected = true
				logger.InfoFields("Event Store is back online, the unsent events that will begin to send.", logrus.Fields{
					"unsent": len(es.queue),
				})
			}
			continue
		}

		// Skip if there's nothing in the queue.
		if len(es.queue) == 0 {
			continue
		}

		// A downward loop for the queue.
		for i := len(es.queue) - 1; i >= 0; i-- {
			e := es.queue[i]
			// Wait a little bit for another event.
			<-time.After(time.Millisecond * 2)
			// create the stream write.
			writer := es.NewStreamWriter(e.stream)
			// Append the event in the stream.
			err := writer.Append(nil, goes.NewEvent("", "", e.data, e.meta))
			if err != nil {
				logger.ErrorFields("Error occurred while pushing the event to the stream of Event Store.", logrus.Fields{
					"stream": e.stream,
					"meta":   e.meta,
				})
				continue
			}
			QueueTotal--
			// Remove the event from the queue since it has been sent.
			es.queue = append(es.queue[:i], es.queue[i+1:]...)
		}
	}
}

// send the event to the specified stream.
func (es *eventstore) send(stream string, data interface{}, meta interface{}) {
	// Prepend the event to the queue to keep the order,
	// because the event queue is a downward loop.
	es.queue = append([]event{event{stream, data, meta}}, es.queue...)

	// Counters.
	SentTotal++
	QueueTotal++

	if !es.isConnected {
		logger.WarningFields("Event will be sent when the Event Store is back online.", logrus.Fields{
			"stream": stream,
			"unsent": len(es.queue),
		})
	}
}

// Send the event to Event Store.
func (es *eventstore) Send(e evt.E) {
	// Fill the metadata if it's empty.
	if e.Metadata == nil {
		e.Metadata = map[string]string{
			"service_version": version.Version,
			"service_id":      sd.ID,
			"service_tags":    sd.Tags,
		}
	}

	go es.send(e.Stream, e.Data, e.Metadata)
}
