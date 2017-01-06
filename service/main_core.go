package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"os/signal"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	consulsd "github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jetbasrawi/go.geteventstore"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

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

type eventListener struct {
	event   string
	body    map[string]interface{}
	meta    map[string]string
	handler func(map[string]interface{}, map[string]string)
}

func createEventStore() *goes.Client {
	client, err := goes.NewClient(nil, os.Getenv("KITSVC_ES_SERVER_URL"))
	if err != nil {
		panic(err)
	}

	client.SetBasicAuth(os.Getenv("KITSVC_ES_USERNAME"), os.Getenv("KITSVC_ES_PASSWORD"))

	return client
}

func setEventSubscription(client *goes.Client, listeners []eventListener) {

	for _, v := range listeners {

		//

		reader := client.NewStreamReader(v.event)

		for reader.Next() {
			if reader.Err() != nil {
				if _, ok := reader.Err().(*goes.ErrNoMoreEvents); ok {
					continue
				} else if _, ok := reader.Err().(*goes.ErrNotFound); ok {
					writer := client.NewStreamWriter(v.event)
					err := writer.Append(nil, goes.NewEvent("", "", map[string]string{}, map[string]string{}))
					if err != nil {
						panic(err)
					}
					continue
				} else if match, _ := regexp.MatchString(".*connection refused.*", reader.Err().Error()); match {
					time.Sleep(time.Second * 2)
					continue
				} else {
					panic(reader.Err())
				}
			}

			err := reader.Scan(&v.body, &v.meta)
			if err != nil {
				panic(err)
			}

			v.handler(v.body, v.meta)
		}
	}
}

//
//
//
//
//

// Model represents the model layer of the service.
type Model struct {
	*sql.DB
}

// createModel creates the model of the service with the database connection.
func createModel(db *sql.DB) Model {
	return Model{db}
}

//
//
//
//
//

type service struct {
	Model
	ES *goes.Client
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
func createService(logger kitlog.Logger, model Model, es *goes.Client) (Service, context.Context) {

	var svc Service
	svc = service{Model: model, ES: es}
	svc = createLoggingMiddleware(logger)(svc)
	svc = createInstruMiddleware()(svc)

	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	setServiceSubscription(serviceHandlers(ctx, options, svc))
	go setEventSubscription(es, eventListeners(svc))

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
	Logger kitlog.Logger
	Service
}

// createLoggingMiddleware creates the logging middleware.
func createLoggingMiddleware(logger kitlog.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return LoggingMiddleware{Logger: logger, Service: next}
	}
}

// createLogger creates the logger with the specified port which tracks the function callers.
func createLogger(port *string) kitlog.Logger {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With("listen", port).With("caller", kitlog.DefaultCaller)

	return logger
}

//
//
//
//
//

// registerService register the service to the service discovery server(consul).
func registerService(logger kitlog.Logger) {
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
func createDatabase(resetDB bool) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
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

	if resetDB {
		databaseUpstream(db)
	}

	return db
}
