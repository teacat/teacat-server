package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	nsq "github.com/bitly/go-nsq"
)

var (
	// ErrEmpty 會在傳入一個空字串時被觸發。
	ErrEmpty = ErrInfo{
		Text:   errors.New("The string is empty."),
		Status: http.StatusBadRequest,
		Code:   "str_empty",
	}
)

// StringService 是基於字串的服務。
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	Test(*nsq.Message)
}

type service struct {
	Message *nsq.Producer
	Model
}

type ServiceMiddleware func(Service) Service

// Uppercase 將傳入的字串轉換為大寫。
func (svc service) Uppercase(s string) (string, error) {

	//c.Message.Publish("new_user", []byte("test"))

	res, err := svc.Model.ToUpper(s)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Count 計算傳入的字串長度。
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
