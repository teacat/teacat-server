package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/TeaMeow/KitSvc/service/store"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	consulsd "github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jetbasrawi/go.geteventstore"
	"github.com/jinzhu/gorm"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

type H map[string]interface{}

// MapString converts type H to map[string]string
func (h H) MapString() map[string]string {
	m := make(map[string]string)

	for k, v := range h {
		m[k] = fmt.Sprint(v)
	}

	return m
}

func main() {

	// Command line flags.
	var (
		listenPort = flag.String("listen", ":"+os.Getenv("KITSVC_PORT"), "HTTP listen address")
		resetDB    = flag.Bool("reinitialize-database", false, "Set true to reinitialize the database, it's useful with the unit testing.")
	)
	flag.Parse()

	// Get the boolean from the pointer by using the reflect package.
	r, _ := strconv.ParseBool(fmt.Sprint(reflect.ValueOf(resetDB).Elem()))

	// Create the logger with the specified listen port.
	logger := createLogger(listenPort)
	// Create the database connection.
	db := createDatabase()
	// Create the store with the database connection.
	s := createStore(r, db)
	// Create the event store.
	es := createEventStore()
	// Create the main service with what it needs.
	_, _, mux := createService(logger, es, s)

	// Log the ports.
	logger.Log("msg", "HTTP", "addr", *listenPort)
	// Start the service and listening to the requests and let the mux router handles every things.
	logger.Log("err", http.ListenAndServe(*listenPort, mux))
}

// InstrumentingMiddleware describes
// a service (as opposed to endpoint) middleware.
//
// The instrumenting middleware instruments
// the performance, time wasted of the operations.
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

// eventListener represents an event listener.
//
// We'll register the event listeners and
// listen for the incoming event
// from the event store.
type eventListener struct {
	event   string
	body    interface{}
	meta    map[string]interface{}
	handler func(interface{}, map[string]interface{})
}

// createEventStore creates a client of the event store.
func createEventStore() *goes.Client {
	client, err := goes.NewClient(nil, os.Getenv("KITSVC_ES_SERVER_URL"))
	if err != nil {
		panic(err)
	}

	client.SetBasicAuth(os.Getenv("KITSVC_ES_USERNAME"), os.Getenv("KITSVC_ES_PASSWORD"))

	return client
}

// setEventSubscription applies the event handlers.
func setEventSubscription(client *goes.Client, logger kitlog.Logger, listeners []eventListener) {

	// The toggle used to detect if the events were all replayed or not.
	played := make(chan kitlog.Logger)
	sent := false

	// Register the service to sd if the events were all replayed.
	go func() {
		// If the events were all replayed.
		if l := <-played; l != nil {
			// then register the service to the service registry.
			registerService(l)
			// End the goroutine
			return
		}
	}()

	// Each of the listener.
	for _, v := range listeners {
		// Create the the stream reader for listening the specified event(stream).
		reader := client.NewStreamReader(v.event)

		// Read the next event.
		for reader.Next() {

			if reader.Err() != nil {
				switch reader.Err().(type) {

				// Continue if there's no more event.
				case *goes.ErrNoMoreEvents:

					// Since there're no more messages can read. We've replayed all the events,
					// and it's time to register the service to the sd because we're ready.
					if !sent {
						// Send the logger to the played channel because we need the logger.
						played <- logger
						// Set the sent toggle as true so we won't send the logger to the channel again.
						sent = true
						// Close the unused channel.
						close(played)
					}

					// When there are no more event in the stream, set LongPoll.
					// The server will wait for 5 seconds in this case or until
					// events become available on the stream.
					reader.LongPoll(5)

				// Create an empty event if the stream hasn't been created.
				case *goes.ErrNotFound:
					writer := client.NewStreamWriter(v.event)

					// Create am empty stream.
					err := writer.Append(nil, goes.NewEvent("", "", map[string]string{}, map[string]string{}))
					if err != nil {
						panic(err)
					}
					continue

				// Sleep for 5 seconds and try again if the EventStore was not connected.
				case *url.Error, *goes.ErrTemporarilyUnavailable:
					fmt.Println("Cannot conenct to the Event Sotre, try again after 5 seconds.")
					<-time.After(time.Duration(5) * time.Second)

				// Bye bye if really error.
				default:
					panic(reader.Err())
				}

			} else {
				// Mapping the event data.
				err := reader.Scan(&v.body, &v.meta)
				if err != nil {
					//panic(err)
				}

				// Skip if the body or the meta is empty (might be the empty one which we created for the new stream).
				//if len(v.body) == 0 {
				//	continue
				//}

				// Call to the event handler.
				v.handler(v.body, v.meta)
			}
		}
	}
}

// createStore creates the store with the database connection,
// also creates/destorys the database tables.
//
// The store is used to interactive with the database,
// we'll modify the database from the store,
// not from the controller.
func createStore(resetDB bool, db *gorm.DB) store.Store {
	s := store.Store{DB: db}

	if resetDB {
		s.Downstream()
		s.Upstream()
	}

	return s
}

// service represents the service concrete.
//
// The service concrete contains
// the operations of the service can do,
// and it also contains the modules that we can use in the operations.
type service struct {
	Store store.Store
	ES    *goes.Client
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

// Err represents the error response.
type Err struct {
	Message error
	Payload interface{}
}

// Error returns the error string from the Err struct.
func (e Err) Error() string {
	return e.Message.Error()
}

// ErrInfo represents the information of the operation error.
type ErrInfo struct {
	Text   error
	Status int
	Code   string
}

// Error returns the error string from the ErrInfo struct.
func (e ErrInfo) Error() string {
	return e.Text.Error()
}

// response represents the response of the endpoints.
type response struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// errorEncoder encodes the error information to an endpoint response.
func errorEncoder(c context.Context, err error, w http.ResponseWriter) {
	var status int
	var code string
	var msg string
	var payload interface{}

	switch err.(type) {

	// The error which creates by the operations of the service.
	case Err:
		status, msg, code, payload =
			err.(Err).Message.(ErrInfo).Status,
			err.(Err).Message.(ErrInfo).Text.Error(),
			err.(Err).Message.(ErrInfo).Code,
			err.(Err).Payload

	// The unknown error, JSON parse error most of the time.
	default:
		status, msg, code, payload =
			http.StatusBadRequest,
			"Cannot parse the JSON content.",
			"error",
			nil
	}

	// Set the JSON content type.
	w.Header().Set("Content-Type", "application/json")
	// Set the status code.
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response{
		Status:  "error",
		Code:    code,
		Message: msg,
		Payload: payload,
	})
}

// encodeResponse encodes the data to a response.
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
func createService(logger kitlog.Logger, es *goes.Client, s store.Store) (Service, context.Context, *mux.Router) {

	// Create the service and the middlewares.
	var svc Service
	svc = service{Store: s, ES: es}
	svc = createLoggingMiddleware(logger)(svc)
	svc = createInstruMiddleware()(svc)

	// The context and the server options.
	ctx := context.Background()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}

	// Set the operation handlers and get the mux router.
	mux := setServiceSubscription(ctx, options, svc, serviceHandlers(ctx, options, svc))
	// Set the event handlers with goroutine so it won't blocked.
	go setEventSubscription(es, logger, eventListeners(svc))

	return svc, ctx, mux
}

// serviceHandler represents a service handler.
type serviceHandler struct {
	method  string
	pattern string
	handler http.Handler
}

// setServiceSubscription sets the service handlers.
func setServiceSubscription(ctx context.Context, opts []httptransport.ServerOption, svc Service, handlers []serviceHandler) *mux.Router {

	// Create a mux router to handle the incoming requests.
	r := mux.NewRouter()

	// The service discovery health check handler.
	consulsdHandler := httptransport.NewServer(ctx, makeServiceDiscoveryEndpoint(svc), decodeServiceDiscoveryRequest, encodeResponse, opts...)
	r.Handle("/sd_health", consulsdHandler).Methods("GET")

	// The metrics handler for prometheus.
	r.Handle("/metrics", stdprometheus.Handler()).Methods("GET")

	// Apply the other custom handlers.
	for _, v := range handlers {
		r.Handle(v.pattern, v.handler).Methods(v.method)
	}

	return r
}

// LoggingMiddleware describes
// a service (as opposed to endpoint) middleware.
//
// The logging middleware logs the data of the incoming request
// and the output data.
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

// registerService registers the service
// to the service discovery server(consul).
//
// So the gateway can call to the service
// without the extra settings or restarting.
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

// sdResponse represents the response of the health check for consul.
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

// createDatabase creates the database connection by using gorm.
// https://github.com/jinzhu/gorm
//
//
//
func createDatabase() *gorm.DB {

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

	return db
}
