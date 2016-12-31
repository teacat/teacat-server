package service

import (
	"encoding/json"
	"net/http"

	nsq "github.com/bitly/go-nsq"
	"github.com/jinzhu/gorm"

	"golang.org/x/net/context"
)

// stringService 概括了字串服務所可用的函式。
type Concrete struct {
	Message *nsq.Producer
	Model
}

type Model struct {
	*gorm.DB
}

// ServiceMiddleware 是處理 StringService 的中介層。
type Middleware func(Service) Service

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
