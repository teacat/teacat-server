package main

import (
	"errors"
	"fmt"
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
	Test(*nsq.Message)
}

// Service operation is just like the controller in the MVC architecture,
// We don't process the data in the controller but decide what model to call,
// then we pass the data to the model.
//
// Create the service operations with the following format:
//     func (svc service)...

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
