package eventstore

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/jetbasrawi/go.geteventstore"
	"github.com/parnurzeal/gorequest"
)

type eventstore struct {
	*goes.Client
	isConnected bool
	queue       []event
}

type event struct {
	stream string
	data   interface{}
	meta   interface{}
}

// newClient creates a new event store client.
func NewClient(url string, esURL string, username string, password string, e *eventutil.Engine, played chan<- bool, ready <-chan bool) *eventstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingStore(esURL); err != nil {
		logrus.Fatalln(err)
	}
	// Create the client to Event Store.
	client, err := goes.NewClient(nil, esURL)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Cannot create the client of Event Store.")
	}
	// Set the username, password
	client.SetBasicAuth(username, password)

	es := &eventstore{client, true, []event{}}

	// Capturing the events when the router was ready in the goroutine.
	go es.capture(url, e, played, ready)

	// Pushing the events in the local queue to Event Store.
	go es.push(esURL)

	return es
}

//
func pingStore(url string) error {
	for i := 0; i < 30; i++ {
		// Ping the Event Store by sending a GET request,
		// and the response status code doesn't matter.
		_, err := http.Get(url)
		if err == nil {
			return nil
		}

		// Waiting for another round if we didn't receive the response.
		logrus.Infof("Cannot connect to Event Store, retry in 3 second.")
		time.Sleep(time.Second * 3)
	}

	return errors.New("Cannot connect to Event Store.")
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
		logrus.Fatalln("Error occurred while sending the event to self router.")
	}
	if resp.StatusCode != 200 {
		logrus.Infoln("The event has been recevied by the router, but the status code wasn't 200.")
	}
}

// isEmptyEvent returns true when the event body is empty.
func isEmptyEvent(json []byte) bool {
	return string(json) == "{}"
}

//
func (es *eventstore) capture(localUrl string, e *eventutil.Engine, played chan<- bool, ready <-chan bool) {
	// Continue if the router was ready.
	<-ready

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
							playedStreams++
							replayedStream = true

							// Show the statistic of the replayed stream.
							logrus.Infof("The `%s` events were all replayed in %.2f seconds, there's a total of %d events were in the stream.",
								l.Stream,
								time.Since(startTime).Seconds(),
								totalEvents,
							)
							// Set the `played` as true once we have replayed all the event streams.
							if playedStreams >= totalStreams {
								played <- true
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
							logrus.Errorln(err)
							logrus.Fatalln("Error occurred while creating an empty stream.")
						}
						continue

						// Sleep for few seconds and try again if the Event Store was not connected.
					case *url.Error, *goes.ErrTemporarilyUnavailable:
						es.isConnected = false
						<-time.After(time.Duration(3) * time.Second)

						// Other errors.
					default:
						logrus.Warningln(r.Err())
						logrus.Warningln("Error occurred while reading the incoming event.")
						continue
					}

					// We received the event, and we're going to make a http request,
					// abd send it to out own Gin router.
				} else {
					// Get the event body.
					json, err := r.EventResponse().Event.Data.(*json.RawMessage).MarshalJSON()
					if err != nil {
						logrus.Warningln(err)
						logrus.Warningf("Cannot parse the event data, the `%s` event has been skipped.", l.Stream)
						continue
					}
					// Skip the empty json event,
					// because that might be the one which we used to create the new stream.
					if isEmptyEvent(json) {
						//logrus.Infof("Received the empty `%s` event.", l.Stream)
						continue
					}

					totalEvents++
					// Send the received data to the router,
					// so we can process the event in the router.
					go sendToRouter(l.Method, localUrl+l.Path, json)

					continue
				}
			}
		}(l)
	}
}

func (es *eventstore) push(esURL string) {
	for {
		// Check the queue every second.
		<-time.After(time.Second * 1)

		//calculate send rateRATE

		// Ping the Event Store to see if it's back online or not.
		if !es.isConnected {
			if err := pingStore(esURL); err == nil {
				es.isConnected = true

				logrus.Infof("Event Store is back online, there are %d unsent events that will begin to send.", len(es.queue))
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

			// create the stream write.
			writer := es.NewStreamWriter(e.stream)
			// Append the event in the stream.
			err := writer.Append(nil, goes.NewEvent("", "", e.data, e.meta))
			if err != nil {
				// ERR HANDLE
				continue
			}

			// Remove the event from the queue since it has been sent.
			es.queue = append(es.queue[:i], es.queue[i+1:]...)
		}
	}
}

func (es *eventstore) send(stream string, data interface{}, meta interface{}) {
	// Prepend to keep the order, because the push queue is a downward loop.
	es.queue = append([]event{event{stream, data, meta}}, es.queue...)
	//
	if !es.isConnected {
		logrus.Warningf("The `%s` event will be sent when the Event Store is back online. (Queue length: %d)", stream, len(es.queue))
	}
}
