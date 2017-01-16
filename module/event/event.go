package event

import (
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/codegangsta/cli"
	"github.com/jetbasrawi/go.geteventstore"
)

// newClient creates a new event store client.
func newClient(c *cli.Context) *goes.Client {
	client, err := goes.NewClient(nil, c.String("es-url"))
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Event store connection failed.")
	}

	client.SetBasicAuth(c.String("es-username"), c.String("es-password"))

	return client
}

func Capture(c *cli.Context, e *eventutil.Engine, played chan<- bool) {

	sent := false

	// Create a new client.
	client := newClient(c)

	// Each of the listener.
	for _, l := range e.Listeners {

		// Create the the stream reader for listening the specified stream.
		reader := client.NewStreamReader(l.Stream)

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
						played <- true
						// Set the sent toggle as true so we won't send the logger to the channel again.
						sent = true
					}

					// When there are no more events in the stream, set LongPoll.
					// The server will wait for 5 seconds in this case or until
					// events become available on the stream.
					reader.LongPoll(5)

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
					logrus.Warningln("Cannot connect to the event store, try again after 5 seconds.")
					<-time.After(time.Duration(5) * time.Second)

					// Bye bye if really error.
				default:
					logrus.Errorln(reader.Err())
					logrus.Fatalln("Error occurred while reading the incoming event.")
				}

				// Receivied the event.
			} else {

			}
		}
	}
}
