package main

import (
	"errors"
	"net/http"

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

// Service represents the operations of the service can do.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	PublishMessage(string)
	ReceiveMessage(*nsq.Message)
	ServiceDiscoveryCheck()
}

// Service operation is just like the controller in the MVC architecture,
// We don't process the data in the controller but decide what model to call,
// then we pass the data to the model.
//
// Create the service operations with the following format:
//     func (svc service)...

// Uppercase converts the string to uppercase.
func (svc service) Uppercase(s string) (string, error) {

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

func (svc service) PublishMessage(s string) {
	svc.Message.Publish("hello_world", []byte(s))
}

func (service) ReceiveMessage(msg *nsq.Message) {
}

func (service) ServiceDiscoveryCheck() {

}
