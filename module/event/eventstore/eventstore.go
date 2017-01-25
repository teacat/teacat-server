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
}

// newClient creates a new event store client.
func NewClient(url string, esUrl string, username string, password string, e *eventutil.Engine, isPlayed chan<- bool, isReady <-chan bool) *eventstore {
	// Ping the Event Store to make sure it's alive.
	if err := pingStore(esUrl); err != nil {
		logrus.Fatalln("Cannot connect to Event Store after retried so many times.")
	}
	// Create the client to Event Store.
	client, err := goes.NewClient(nil, esUrl)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Cannot create the client of Event Store.")
	}

	// Set the username, password
	client.SetBasicAuth(username, password)

	// Capturing the events when the router was ready in the goroutine.
	go capture(url, client, e, isPlayed, isReady)

	return &eventstore{client}
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

//
func capture(localUrl string, client *goes.Client, e *eventutil.Engine, played chan<- bool, ready <-chan bool) {
	// Continue if the router was ready.
	<-ready

	// The `All events were played`` state.
	playedCount := 0
	totalCount := len(e.Listeners)

	// Each of the listener.
	for _, l := range e.Listeners {
		// Create the the stream reader for listening the specified stream.
		r := client.NewStreamReader(l.Stream)

		go func(l eventutil.Listener) {
			// The `sent` toggle used to make sure if we have sent the `played` indicator or not.
			sent := false
			//
			startTime := time.Now()
			//
			totalEvents := 0

			// Read the next event.
			for r.Next() {
				// Error occurred.
				if r.Err() != nil {
					switch r.Err().(type) {

					// Continue if there's no more event.
					case *goes.ErrNoMoreEvents:
						// Since there're no more messages can read. We've replayed all the events,
						// and it's time to register the service to the sd because we're ready.
						if !sent {
							// Send the logger to the played channel because we need the logger.
							playedCount++
							// Set the sent toggle as true so we won't send the logger to the channel again.
							sent = true

							logrus.Infof("The `%s` events were all replayed in %.2f seconds, there're total %d events in the stream.", l.Stream, time.Since(startTime).Seconds(), totalEvents)

							if playedCount >= totalCount {
								played <- true
							}
						}

						// When there are no more events in the stream, set LongPoll.
						// The server will wait for 15 seconds in this case or until
						// events become available on the stream.
						r.LongPoll(15)

					// Create an empty event if the stream hasn't been created.
					case *goes.ErrNotFound:
						writer := client.NewStreamWriter(l.Stream)

						// Create am empty stream.
						err := writer.Append(nil, goes.NewEvent("", "", map[string]string{}, map[string]string{}))
						if err != nil {
							logrus.Errorln(err)
							logrus.Fatalln("Error occurred while creating an empty stream.")
						}
						continue

						// Sleep for 5 seconds and try again if the EventStore was not connected.
					case *url.Error, *goes.ErrTemporarilyUnavailable:
						logrus.Warningln("Cannot connect to the event store, try again after 3 seconds.")
						<-time.After(time.Duration(3) * time.Second)

						// Bye bye if really error.
					default:
						logrus.Warningln(r.Err())
						logrus.Warningln("Error occurred while reading the incoming event.")
						continue
					}

					// We received the event, and we're going to make a http request,
					// send it to out own Gin router.
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

					// Send the received data to self-router,
					// so we can process it in the router.
					go sendToRouter(l.Method, localUrl+l.Path, json)
					continue
				}
			}
		}(l)
	}
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

/*
func (es *eventstore) CreateUser() {

}*/
