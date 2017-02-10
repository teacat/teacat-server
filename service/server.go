package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeaMeow/KitSvc/module/sd"
	"github.com/TeaMeow/KitSvc/router"
	"github.com/TeaMeow/KitSvc/router/middleware"
	"github.com/TeaMeow/KitSvc/shared/eventutil"
	"github.com/TeaMeow/KitSvc/shared/mqutil"
	"github.com/TeaMeow/KitSvc/shared/wsutil"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

var serverFlags = []cli.Flag{
	// Common flags.
	cli.StringFlag{
		EnvVar: "KITSVC_NAME",
		Name:   "name",
		Usage:  "the name of the service, exposed for service discovery.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_URL",
		Name:   "url",
		Usage:  "the url of the service.",
		Value:  "http://127.0.0.1:8080",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_ADDR",
		Name:   "addr",
		Usage:  "the address of the service (with the port).",
		Value:  "127.0.0.1:8080",
	},
	cli.IntFlag{
		EnvVar: "KITSVC_PORT",
		Name:   "port",
		Usage:  "the port of the service.",
		Value:  8080,
	},
	cli.StringFlag{
		EnvVar: "KITSVC_USAGE",
		Name:   "usage",
		Usage:  "the usage of the service, exposed for service discovery.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_JWT_SECRET",
		Name:   "jwt-secret",
		Usage:  "the secert used to encode the json web token.",
	},
	//cli.StringFlag{
	//	EnvVar: "KITSVC_VERSION",
	//	Name:   "version",
	//	Usage:  "the version of the service.",
	//},

	// Database flags.
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_NAME",
		Name:   "database-name",
		Usage:  "the name of the database.",
		Value:  "service",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_HOST",
		Name:   "database-host",
		Usage:  "the host of the database (with the port).",
		Value:  "127.0.0.1:3306",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_USER",
		Name:   "database-user",
		Usage:  "the user of the database.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_PASSWORD",
		Name:   "database-password",
		Usage:  "the password of the database.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_CHARSET",
		Name:   "database-charset",
		Usage:  "the charset of the database.",
		Value:  "utf8",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_DATABASE_LOC",
		Name:   "database-loc",
		Usage:  "the timezone of the database.",
		Value:  "Local",
	},
	cli.BoolFlag{
		EnvVar: "KITSVC_DATABASE_PARSE_TIME",
		Name:   "database-parse_time",
		Usage:  "parse the time.",
	},

	// NSQ flags.
	cli.StringFlag{
		EnvVar: "KITSVC_NSQ_PRODUCER",
		Name:   "nsq-producer",
		Usage:  "the address of the TCP NSQ producer (with the port).",
		Value:  "127.0.0.1:4150",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_NSQ_PRODUCER_HTTP",
		Name:   "nsq-producer-http",
		Usage:  "the address of the HTTP NSQ producer (with the port).",
		Value:  "127.0.0.1:4151",
	},
	cli.StringSliceFlag{
		EnvVar: "KITSVC_NSQ_LOOKUPDS",
		Name:   "nsq-lookupds",
		Usage:  "the address of the NSQ lookupds (with the port).",
		Value: &cli.StringSlice{
			"127.0.0.1:4161",
		},
	},

	// Event store flags.
	cli.StringFlag{
		EnvVar: "KITSVC_ES_SERVER_URL",
		Name:   "es-url",
		Usage:  "the url of the event store server.",
		Value:  "http://127.0.0.1:2113",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_ES_USERNAME",
		Name:   "es-username",
		Usage:  "the username of the event store.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_ES_PASSWORD",
		Name:   "es-password",
		Usage:  "the password of the event store.",
	},

	// Prometheus flags.
	cli.StringFlag{
		EnvVar: "KITSVC_PROMETHEUS_NAMESPACE",
		Name:   "prometheus-namespace",
		Usage:  "the prometheus namespace.",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_PROMETHEUS_SUBSYSTEM",
		Name:   "prometheus-subsystem",
		Usage:  "the subsystem of the promethues.",
	},

	// Consul flags.
	cli.StringFlag{
		EnvVar: "KITSVC_CONSUL_CHECK_INTERVAL",
		Name:   "consul-check_interval",
		Usage:  "the interval of consul health check.",
		Value:  "10s",
	},
	cli.StringFlag{
		EnvVar: "KITSVC_CONSUL_CHECK_TIMEOUT",
		Name:   "consul-check_timeout",
		Usage:  "the timeout of consul health check.",
		Value:  "1s",
	},
	cli.StringSliceFlag{
		EnvVar: "KITSVC_CONSUL_TAGS",
		Name:   "consul-tags",
		Usage:  "the service tags for consul.",
		Value: &cli.StringSlice{
			"micro",
		},
	},
}

func server(c *cli.Context, serverReady chan<- bool) error {
	// The ready, event played states.
	isReady := make(chan bool, 2)
	isPlayed := make(chan bool)

	// Create the Gin engine.
	gin := gin.New()
	// Create the event handler struct.
	event := eventutil.New(gin)
	//
	ws := wsutil.New(gin)
	//
	mq := mqutil.New(gin)

	// Routes.
	router.Load(
		gin,
		event,
		ws,
		mq,
		middleware.Config(c),
		middleware.Store(c),
		middleware.Logging(),
		middleware.Event(c, event, isPlayed, isReady),
		middleware.MQ(c, mq, isReady),
		middleware.Metrics(),
	)

	// And register the service to the service registry when the events were replayed in the goroutine.
	go sd.Wait(c, isPlayed)

	// We only do those things when the router is ready to use.
	go func() {
		// To check the router is good to go,
		// we ping the server by sending the GET request to the router.
		if err := pingServer(c); err != nil {
			logrus.Fatalln("The router has no response, or it might took too long to start up.")
		}

		logrus.Infoln("The router has been deployed successfully.")
		// Send `true` to the `isReady` channel if the router is ready to use.
		isReady <- true
		isReady <- true
		serverReady <- true
	}()

	// Start to listening the incoming requests.
	return http.ListenAndServe(c.String("addr"), gin)
}

// pingServer pings the http server to make sure the router is currently working.
func pingServer(c *cli.Context) error {
	for i := 0; i < 30; i++ {

		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(c.String("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Waiting for another round if we didn't receive the 200 status code by the ping request.
		logrus.Infof("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}
