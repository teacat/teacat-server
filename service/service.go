package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	nsq "github.com/bitly/go-nsq"
)

// Error codes returned by failures
var (

	// ErrEmpty will returned if the string is empty.
	ErrEmpty = ErrInfo{
		Text:   errors.New("The string is empty."),
		Status: http.StatusBadRequest,
		Code:   "str_empty",
	}
)

// Service represents the operations of the servicer can do.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	Test(*nsq.Message)
}

type service struct {
	Message *nsq.Producer
	Model
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

// Uppercase converts the string to uppercase.
func (svc service) Uppercase(s string) (string, error) {

	//svc.Message.Publish("new_user", []byte("test"))

	res, err := svc.Model.ToUpper(s)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Count counts the length of the string.
func (svc service) Count(s string) int {
	return svc.Model.Count(s)
}

func (service) Test(msg *nsq.Message) {
	fmt.Println(msg)
}

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

	status, msg, code, payload :=
		err.(Err).Message.(ErrInfo).Status,
		err.(Err).Message.(ErrInfo).Text.Error(),
		err.(Err).Message.(ErrInfo).Code,
		err.(Err).Payload

	if status == 0 {
		status = http.StatusInternalServerError
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
