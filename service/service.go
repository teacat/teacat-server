package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TeaMeow/KitSvc/service/event"
)

// Error codes returned by failures
var (
	// ErrEmpty is returned if the string is empty.
	ErrEmpty = ErrInfo{
		Text:   errors.New("The string is empty."),
		Status: http.StatusBadRequest,
		Code:   "string_empty",
	}
	// ErrEvent is returned if the error occurred while we're publishing the event.
	ErrEvent = ErrInfo{
		Text:   errors.New("Error occurred while publishing the event."),
		Status: http.StatusInternalServerError,
		Code:   "event_error",
	}
)

// Service represents the operations of the service can do.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) int
	CatchUppercase(interface{}, map[string]interface{})
}

// Count counts the length of the string.
func (svc service) Count(s string) int {
	return len(s)
}

// Uppercase converts the string to uppercase.
func (svc service) Uppercase(s string) (string, error) {
	if s == "" {
		return "", Err{Message: ErrEmpty}
	}

	// Converts the string to uppercase.
	u := strings.ToUpper(s)

	// Create the string record in the database.
	svc.Store.CreateString(s, u)

	//
	option := event.Option{Client: svc.ES, Stream: "uppercase", Meta: nil}
	data := event.String{Input: s, Output: u}

	if err := data.Send(option); err != nil {
		return "", Err{Message: ErrEvent}
	}

	// Get the last record from the database.
	fmt.Println("Last record:")
	fmt.Println(svc.Store.GetLastString())

	return u, nil
}

// CatchUppercase catches the uppercase event, and print the event data.
func (svc service) CatchUppercase(body interface{}, meta map[string]interface{}) {
	b := body.(*event.String)

	fmt.Println("Event received:")
	fmt.Println("Input:" + b.Input + " | Output:" + b.Output)
}
