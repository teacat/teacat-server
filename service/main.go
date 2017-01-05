package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"

	nsq "github.com/bitly/go-nsq"
	"github.com/go-kit/kit/log"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	consulsd "github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

func main() {

	// Command line flags.
	var (
		listenPort = flag.String("listen", ":"+os.Getenv("KITSVC_PORT"), "HTTP listen address")
		resetDB    = flag.Bool("reinitialize-database", false, "Set true to reinitialize the database, it's useful with the unit testing. The database will backed up before the database was reinitialized.")
	)
	flag.Parse()

	// Create the logger with the specified listen port.
	logger := createLogger(listenPort)
	// Create the database connection.
	db := createDatabase(resetDB)
	// Create the model with the database connection.
	model := createModel(db)
	// Create the messaging service with the logger.
	msg := createMessage()

	// Create the main service with what it needs.
	createService(logger, msg, model)
	// Register the service to the service registry.
	registerService(logger)

	// Log the ports.
	logger.Log("msg", "HTTP", "addr", *listenPort)
	// Start the service and listening to the requests.
	logger.Log("err", http.ListenAndServe(*listenPort, nil))
}

//
//
//
//
//

// InstrumentingMiddleware represents a middleware of the instrumenting.
type InstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	Service
}

// createInstruMiddleware creates the instrumenting middleware.
func createInstruMiddleware() ServiceMiddleware {

	fieldKeys := []string{"method", "error"}

	// Number of requests received.
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: os.Getenv("KITSVC_PROMETHEUS_NAMESPACE"),
		Subsystem: os.Getenv("KITSVC_PROMETHEUS_SUBSYSTEM"),
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	// Total duration of requests in microseconds.
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: os.Getenv("KITSVC_PROMETHEUS_NAMESPACE"),
		Subsystem: os.Getenv("KITSVC_PROMETHEUS_SUBSYSTEM"),
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// The result of each count method.
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: os.Getenv("KITSVC_PROMETHEUS_NAMESPACE"),
		Subsystem: os.Getenv("KITSVC_PROMETHEUS_SUBSYSTEM"),
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	return func(next Service) Service {
		return InstrumentingMiddleware{requestCount, requestLatency, countResult, next}
	}
}

//
//
//
//
//

type messageHandlerFunc func(*nsq.Message)

type messageHandler struct {
	topic   string
	channel string
	handler messageHandlerFunc
}

func createMessage() *nsq.Producer {
	prod, err := nsq.NewProducer(os.Getenv("KITSVC_NSQ_PRODUCER"), nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	return prod
}

func messageSubscribe(topic string, ch string, fn messageHandlerFunc) {

	q, err := nsq.NewConsumer(topic, ch, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	q.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		fn(msg)

		return nil
	}))

	if err := q.ConnectToNSQLookupds(strings.Split(os.Getenv("KITSVC_NSQ_LOOKUPS"), ",")); err != nil {
		panic(err)
	}
}

func setMessageSubscription(handlers []messageHandler) {

	for _, v := range handlers {
		// Create the topic
		cmd := exec.Command("curl", "-X", "POST", "http://127.0.0.1:4151/topic/create?topic="+v.topic)
		cmd.Start()
		cmd.Wait()

		// Subscribe to the topic
		messageSubscribe(v.topic, v.channel, v.handler)
	}
}

//
//
//
//
//

// Model represents the model layer of the service.
type Model struct {
	DB *gorm.DB
}

// createModel creates the model of the service with the database connection.
func createModel(db *gorm.DB) Model {
	return Model{db}
}

//
//
//
//
//

type service struct {
	Message *nsq.Producer
	Model
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

type Err struct {
	Message error
	Payload interface{}
}

func (e Err) Error() string {
	return e.Message.Error()
}

type ErrInfo struct {
	Text   error
	Status int
	Code   string
}

func (e ErrInfo) Error() string {
	return e.Text.Error()
}

type response struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func errorEncoder(c context.Context, err error, w http.ResponseWriter) {
	var status int
	var code string
	var msg string
	var payload interface{}

	switch err.(type) {
	case Err:
		status, msg, code, payload =
			err.(Err).Message.(ErrInfo).Status,
			err.(Err).Message.(ErrInfo).Text.Error(),
			err.(Err).Message.(ErrInfo).Code,
			err.(Err).Payload

	default:
		status, msg, code, payload =
			http.StatusBadRequest,
			"Cannot parse the JSON content.",
			"error",
			nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response{
		Status:  "error",
		Code:    code,
		Message: msg,
		Payload: payload,
	})
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response{
		Status:  "success",
		Code:    "success",
		Message: "",
		Payload: resp,
	})
}

// createService creates the main service by setting the handlers and preparing the middlewares.
func createService(logger kitlog.Logger, msg *nsq.Producer, model Model) (Service, context.Context) {

	var svc Service
	svc = service{Message: msg, Model: model}
	svc = createLoggingMiddleware(logger)(svc)
	svc = createInstruMiddleware()(svc)

	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	setServiceSubscription(serviceHandlers(ctx, options, svc))
	setMessageSubscription(messageHandlers(svc))

	return svc, ctx
}

type serviceHandler struct {
	pattern string
	handler http.Handler
}

func setServiceSubscription(handlers []serviceHandler) {
	for _, v := range handlers {
		http.Handle(v.pattern, v.handler)
	}
}

//
//
//
//
//

// LoggingMiddleware represents a middleware of the logger.
type LoggingMiddleware struct {
	Logger log.Logger
	Service
}

// createLoggingMiddleware creates the logging middleware.
func createLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return LoggingMiddleware{Logger: logger, Service: next}
	}
}

// createLogger creates the logger with the specified port which tracks the function callers.
func createLogger(port *string) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", port).With("caller", log.DefaultCaller)

	return logger
}

//
//
//
//
//

// registerService register the service to the service discovery server(consul).
func registerService(logger log.Logger) {
	p, _ := strconv.Atoi(os.Getenv("KITSVC_PORT"))

	info := consulapi.AgentServiceRegistration{
		Name: os.Getenv("KITSVC_NAME"),
		Port: p,
		Tags: strings.Split(os.Getenv("KITSVC_CONSUL_TAGS"), ","),
		Check: &consulapi.AgentServiceCheck{
			HTTP:     os.Getenv("KITSVC_URL") + "/sd_health",
			Interval: os.Getenv("KITSVC_CONSUL_CHECK_INTERVAL"),
			Timeout:  os.Getenv("KITSVC_CONSUL_CHECK_TIMEOUT"),
		},
	}

	apiConfig := consulapi.DefaultConfig()
	apiClient, _ := consulapi.NewClient(apiConfig)
	client := consulsd.NewClient(apiClient)
	reg := consulsd.NewRegistrar(client, &info, logger)

	// Deregister the service when exiting the program.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			reg.Deregister()
			os.Exit(1)
		}
	}()

	// Register the service.
	reg.Register()
}

type sdResponse struct {
	P string `json:"pong"`
}

func makeServiceDiscoveryEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return sdResponse{"pong"}, nil
	}
}

func decodeServiceDiscoveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

//
//
//
//
//

// createDatabase creates the database connection.
func createDatabase(resetDB *bool) *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("KITSVC_DATABASE_USER"),
		os.Getenv("KITSVC_DATABASE_PASSWORD"),
		os.Getenv("KITSVC_DATABASE_HOST"),
		os.Getenv("KITSVC_DATABASE_NAME"),
		os.Getenv("KITSVC_DATABASE_CHARSET"),
		os.Getenv("KITSVC_DATABASE_PARSE_TIME"),
		os.Getenv("KITSVC_DATABASE_LOC"),
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}
