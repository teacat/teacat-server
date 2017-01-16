package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/event"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router"
	"github.com/TeaMeow/KitSvc/router/middleware"
	"github.com/codegangsta/cli"
)

func server(c *cli.Context) {

	// The gin router and the event handlers.
	serviceHandler, eventHandler := router.Load(
		middleware.Store(c),
		middleware.Logging(),
	)

	// Start to listening the incoming requests.
	go http.ListenAndServe(c.String("addr"), serviceHandler)

	// Wait for the server is ready to serve.
	if err := pingServer(c); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("The router has no response, or it might took too long to startup.")
	}

	evtPlayed := make(chan bool)

	// Then we capturing the events.
	go event.Capture(c, eventHandler, evtPlayed)

	// And we register the service to the service registry.
	go sd.Wait(c, evtPlayed)
}

// pingServer pings the http server to make sure the router is currently working.
func pingServer(c *cli.Context) (err error) {
	for i := 0; i < 30; i++ {

		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(c.String("url") + "/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Wait for another round if we didn't receive the 200 status code by the ping request.
		logrus.Infof("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}
