package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	CatchEvent(map[string]interface{}, map[string]string)
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

	u := strings.ToUpper(s)

	svc.Store.CreateString(s, u)

	fmt.Println(svc.Store.GetLastString())

	return u, nil
}

func (svc service) CatchEvent(body map[string]interface{}, meta map[string]string) {
	fmt.Println(body["test"])
}
