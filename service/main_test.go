package main

import (
	"bytes"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	testSvc    Service
	testCtx    context.Context
	testLogger kitlog.Logger
)

func init() {
	listenPort := ":8080"

	// Create the logger with the specified listen port.
	testLogger := createLogger(&listenPort)
	// Create the database connection.
	db := createDatabase()
	// Create the store with the database connection.
	s := createStore(false, db)
	// Create the event store.
	es := createEventStore()

	// Create the main service with what it needs.
	svcLocal, ctxLocal, muxLocal := createService(testLogger, es, s)
	testSvc, testCtx = svcLocal, ctxLocal

	// Sleep a little bit till we registered to the sd.
	time.Sleep(time.Second * 1)
	// Start the service and listening to the requests and let the mux router handles every things.
	go http.ListenAndServe(listenPort, muxLocal)
}

// testFunction sends the JSON request, and received the JSON response.
func testFunction(method string, pattern string, body string) (string, error) {

	// Converts the byte array to json
	b := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, os.Getenv("KITSVC_URL")+pattern, b)
	req.Header.Set("Content-Type", "application/json")

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Converts the JSON response to string.
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}
	// Removes the newline symbol in the end of the string.
	s := strings.TrimSpace(buf.String())

	return s, err
}

// TestES tests the event store.
func TestES(t *testing.T) {
	es := createEventStore()
	randomEvent := rand.New(rand.NewSource(time.Now().UnixNano()))

	go setEventSubscription(es, testLogger, []eventListener{
		{
			event:   strconv.Itoa(randomEvent.Intn(999999)),
			body:    make(H),
			meta:    make(H),
			handler: testSvc.CatchUppercase,
		}})
}

// TestServiceDiscoveryMessage tests the health check handler for service discovery server.
func TestServiceDiscoveryMessage(t *testing.T) {
	{
		body := ``
		expected := `{"status":"success","code":"success","message":"","payload":{"pong":"pong"}}`
		resp, _ := testFunction("GET", "/sd_health", body)

		assert.Equal(t, expected, resp, "Cannot ping to the service discovery health check function.")
	}
}

// TestDatabase tests the database connection.
func TestDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Database didn't panic when the host is incorrect.")
		}
	}()

	// Test the database migration.
	db := createDatabase()
	createStore(true, db)

	// Test the incorrect host.
	os.Setenv("KITSVC_DATABASE_HOST", "xxxxxx")
	createDatabase()
}
