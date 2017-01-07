package main

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	svc Service
	ctx context.Context
)

func init() {
	listenPort := ":8080"
	resetDB := false

	// Create the logger with the specified listen port.
	logger := createLogger(&listenPort)
	// Create the database connection.
	db := createDatabase()
	//
	s := createStore(resetDB, db)
	//
	es := createEventStore()

	// Create the main service with what it needs.

	// Create the main service with what it needs.
	svcLocal, ctxLocal, muxLocal := createService(logger, es, s)
	svc, ctx = svcLocal, ctxLocal

	registerService(logger)
	time.Sleep(time.Second * 1)
	http.Handle("/", muxLocal)
	go http.ListenAndServe(listenPort, nil)
}

func testFunction(method string, pattern string, body string) (string, error) {

	b := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, os.Getenv("KITSVC_URL")+pattern, b)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		panic(err)
	}
	s := strings.TrimSpace(buf.String())

	return s, err
}

func TestES(t *testing.T) {
	es := createEventStore()
	go setEventSubscription(es, []eventListener{
		{
			event:   "xxxxxxxxxxx",
			body:    make(map[string]interface{}),
			meta:    make(map[string]string),
			handler: svc.CatchEvent,
		}})
}

func TestServiceDiscoveryMessage(t *testing.T) {
	{
		body := ``
		expected := `{"status":"success","code":"success","message":"","payload":{"pong":"pong"}}`
		resp, _ := testFunction("GET", "/sd_health", body)

		assert.Equal(t, expected, resp, "Cannot ping to the service discovery health check function.")
	}
}

func TestDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Database didn't panic when the host is incorrect.")
		}
	}()

	db := createDatabase()
	createStore(true, db)

	os.Setenv("KITSVC_DATABASE_HOST", "xxxxxx")
	createDatabase()

}
