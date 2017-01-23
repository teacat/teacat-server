package eventstore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/jetbasrawi/go.geteventstore"
)

type eventstore struct {
	*goes.Client
}

// newClient creates a new event store client.
func NewClient(url string, esUrl string, username string, password string, e *eventutil.Engine, isPlayed chan<- bool, isReady <-chan bool) *eventstore {
	client, err := goes.NewClient(nil, esUrl)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Event store connection failed.")
	}
	// Set the username, password
	client.SetBasicAuth(username, password)
	// Capturing the events when the router was ready in the goroutine.
	go capture(url, client, e, isPlayed, isReady)

	return &eventstore{client}
}

//
func capture(localUrl string, client *goes.Client, e *eventutil.Engine, played chan<- bool, ready <-chan bool) {
	// Continue if the router was ready.
	<-ready

	// The `All events were played`` state.
	allPlayed := make(chan bool, len(e.Listeners))

	// Keep watching if all the events were replayed or not.
	go func() {
		for {
			if len(allPlayed) == cap(allPlayed) {
				played <- true
				return
			}
		}
	}()

	// Each of the listener.
	for _, l := range e.Listeners {
		// Create the the stream reader for listening the specified stream.
		reader := client.NewStreamReader(l.Stream)

		go func(l eventutil.Listener) {

			// The `sent` toggle used to make sure if we have sent the `played` indicator or not.
			sent := false

			// Read the next event.
			for reader.Next() {

				// Error occurred.
				if reader.Err() != nil {
					switch reader.Err().(type) {

					// Continue if there's no more event.
					case *goes.ErrNoMoreEvents:

						// Since there're no more messages can read. We've replayed all the events,
						// and it's time to register the service to the sd because we're ready.
						if !sent {
							// Send the logger to the played channel because we need the logger.
							allPlayed <- true
							// Set the sent toggle as true so we won't send the logger to the channel again.
							sent = true
						}

						// When there are no more events in the stream, set LongPoll.
						// The server will wait for 15 seconds in this case or until
						// events become available on the stream.
						reader.LongPoll(15)

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
						logrus.Warningln(reader.Err())
						logrus.Warningln("Error occurred while reading the incoming event.")
						continue
					}

					// We received the event, and we're going to make a http request,
					// send it to out own Gin router.
				} else {
					// Get the event body.
					json, _ := reader.EventResponse().Event.Data.(*json.RawMessage).MarshalJSON()

					// Skip the empty json event,
					// because that might be the one which we used to create the new stream.
					if isEmptyEvent(json) {
						logrus.Infof("Received the empty event from the `%s`.", l.Stream)
						continue
					}

					// Send the received data to self-router,
					// so we can process it in the router.
					sendToRouter(l.Method, localUrl+l.Path, json)
				}
			}
		}(l)
	}
}

// sendToRouter sends the received event data to self router.
func sendToRouter(method string, path string, json []byte) {
	// Prepare to send the event data to the gin router.
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	// Send the request via the HTTP client.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Error occurred while sending the event to self router.")
	}
	defer resp.Body.Close()

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
