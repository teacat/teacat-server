package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/endpoint"

	"golang.org/x/net/context"
)

var (
	svc Service
	ctx context.Context
)

func init() {
	// Load the configurations.
	loadConfig("../")

	listenPort := ":8080"
	resetDB := false

	// Create the logger with the specified listen port.
	logger := createLogger(&listenPort)
	// Create the database connection.
	db := createDatabase(&resetDB)
	// Create the model with the database connection.
	model := createModel(db)
	// Create the messaging service with the logger.
	msg := createMessage(logger)

	// Create the main service with what it needs.
	svc, ctx = createService(logger, msg, model)

	go http.ListenAndServe(listenPort, nil)
}

type testEndpoint struct {
	body     string
	decoder  func(_ context.Context, r *http.Request) (interface{}, error)
	endpoint func(svc Service) endpoint.Endpoint
}

func createTestEndpoint(e testEndpoint) (interface{}, error) {
	httpRequest := httptest.NewRequest("POST", "http://localhost/", bytes.NewReader([]byte(e.body)))
	req, _ := e.decoder(ctx, httpRequest)

	return e.endpoint(svc)(ctx, req)
}

func TestUppercase(t *testing.T) {

	resp, _ := http.Post("http://localhost:8080/uppercase", "application/json", bytes.NewReader([]byte(`{"s": "test"}`)))

	//json.NewDecoder(resp.Body).Decode(&yourStuff)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()
	fmt.Println(s)
	//fmt.Printf("%+v\n", resp)

	/*res := uppercaseResponse{V: "TEST"}
	emptyRes := uppercaseResponse{V: ""}

	if r, _ := createTestEndpoint(testEndpoint{
		body:     `{"s": "test"}`,
		decoder:  decodeUppercaseRequest,
		endpoint: makeUppercaseEndpoint,
	}); r != res {
		t.Error("Uppercase() cannot convert the string to uppercase.")
	}

	if _, err := createTestEndpoint(testEndpoint{
		body:     `{"s": ""}`,
		decoder:  decodeUppercaseRequest,
		endpoint: makeUppercaseEndpoint,
	}); err == nil {
		t.Error("Uppercase() cannot tell the error.")
	}

	if r, _ := createTestEndpoint(testEndpoint{
		body:     `{"s": ""}`,
		decoder:  decodeUppercaseRequest,
		endpoint: makeUppercaseEndpoint,
	}); r != emptyRes {
		t.Error("Uppercase() didn't return the empty string when the error occurred.")
	}*/

	/*if r, _ := svc.Uppercase("test"); r != "TEST" {
		t.Error("Uppercase() cannot convert the string to uppercase.")
	}
	if _, err := svc.Uppercase(""); err == nil {
		t.Error("Uppercase() cannot tell the error.")
	}
	if r, _ := svc.Uppercase(""); r != "" {
		t.Error("Uppercase() didn't return the empty string when the error occurred.")
	}*/
}

func TestCount(t *testing.T) {
	if r := svc.Count("test"); r != 4 {
		t.Error("Count() cannot count the length of the string.")
	}
}
