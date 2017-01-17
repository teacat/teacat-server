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

func server(c *cli.Context) error {

	// The gin router and the event handlers.
	serviceHandler, eventHandler := router.Load(
		middleware.Store(c),
		middleware.Logging(),
		//middleware.Instrumenting(),
	)

	isReady := make(chan bool)
	isPlayed := make(chan bool)

	// Capturing the events when the router was ready in the goroutine.
	go event.Capture(c, eventHandler, isPlayed, isReady)

	// And register the service to the service registry when the events were replayed in the goroutine.
	go sd.Wait(c, isPlayed)

	// We only do those things when the router is ready to use.
	go func() {
		// To check the router is good to go,
		// we ping the server by sending the GET request to the router.
		if err := pingServer(c); err != nil {
			logrus.Errorln(err)
			logrus.Fatalln("The router has no response, or it might took too long to startup.")
		}

		logrus.Infoln("The router has been deployed successfully.")
		// Send `true` to the `isReady` channel if the router is ready to use.
		isReady <- true
	}()

	// Start to listening the incoming requests.
	return http.ListenAndServe(
		c.String("addr"),
		serviceHandler,
	)
}

// pingServer pings the http server to make sure the router is currently working.
func pingServer(c *cli.Context) (err error) {
	for i := 0; i < 30; i++ {

		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(c.String("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Wait for another round if we didn't receive the 200 status code by the ping request.
		logrus.Infof("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}
