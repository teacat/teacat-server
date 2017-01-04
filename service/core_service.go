package main

import (
	"encoding/json"
	"net/http"

	nsq "github.com/bitly/go-nsq"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

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
