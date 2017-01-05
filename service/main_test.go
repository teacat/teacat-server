package main

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

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

	os.Setenv("KITSVC_NAME", "StringService")
	os.Setenv("KITSVC_URL", "http://127.0.0.1:8080")
	os.Setenv("KITSVC_ADDR", "127.0.0.1:8080")
	os.Setenv("KITSVC_PORT", "8080")
	os.Setenv("KITSVC_USAGE", "Operations about the string.")
	os.Setenv("KITSVC_VERSION", "0.0.1")
	os.Setenv("KITSVC_DATABASE_NAME", "service")
	os.Setenv("KITSVC_DATABASE_HOST", "127.0.0.1:3306")
	os.Setenv("KITSVC_DATABASE_USER", "root")
	os.Setenv("KITSVC_DATABASE_PASSWORD", "root")
	os.Setenv("KITSVC_DATABASE_CHARSET", "utf8")
	os.Setenv("KITSVC_DATABASE_LOC", "Local")
	os.Setenv("KITSVC_DATABASE_PARSE_TIME", "true")
	os.Setenv("KITSVC_NSQ_PRODUCER", "127.0.0.1:4150")
	os.Setenv("KITSVC_NSQ_LOOKUPS", "127.0.0.1:4161")
	os.Setenv("KITSVC_PROMETHEUS_NAMESPACE", "my_group")
	os.Setenv("KITSVC_PROMETHEUS_SUBSYSTEM", "string_service")
	os.Setenv("KITSVC_CONSUL_CHECK_INTERVAL", "10s")
	os.Setenv("KITSVC_CONSUL_CHECK_TIMEOUT", "1s")
	os.Setenv("KITSVC_CONSUL_TAGS", "string,micro")

	// Create the logger with the specified listen port.
	logger := createLogger(&listenPort)
	// Create the database connection.
	db := createDatabase(&resetDB)
	// Create the model with the database connection.
	model := createModel(db)

	// Create the main service with what it needs.
	svc, ctx = createService(logger, model)

	registerService(logger)

	go http.ListenAndServe(listenPort, nil)
}

func testFunction(pattern string, body string) (string, error) {

	b := bytes.NewReader([]byte(body))

	resp, err := http.Post(os.Getenv("KITSVC_URL")+pattern, "application/json", b)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := strings.TrimSpace(buf.String())

	return s, err
}

func TestUppercase(t *testing.T) {
	{
		body := `{"s": "test"}`
		expected := `{"status":"success","code":"success","message":"","payload":{"v":"TEST"}}`
		resp, _ := testFunction("/uppercase", body)

		assert.Equal(t, resp, expected, "Uppercase() cannot convert the string to uppercase.")
	}
	{
		body := `{"s":`
		expected := `{"status":"error","code":"error","message":"Cannot parse the JSON content.","payload":null}`
		resp, _ := testFunction("/uppercase", body)

		assert.Equal(t, resp, expected, "Uppercase() cannot tell when the parse error occurred.")
	}
	{
		body := `{"s": ""}`
		expected := `{"status":"error","code":"str_empty","message":"The string is empty.","payload":null}`
		resp, _ := testFunction("/uppercase", body)

		assert.Equal(t, resp, expected, "Uppercase() cannot tell the error.")
	}
}

func TestCount(t *testing.T) {
	{
		body := `{"s": "test"}`
		expected := `{"status":"success","code":"success","message":"","payload":{"v":4}}`
		resp, _ := testFunction("/count", body)

		assert.Equal(t, resp, expected, "Count() cannot count the length of the string.")
	}
	{
		body := `{"s":`
		expected := `{"status":"error","code":"error","message":"Cannot parse the JSON content.","payload":null}`
		resp, _ := testFunction("/count", body)

		assert.Equal(t, resp, expected, "Count() cannot tell when the parse error occurred.")
	}
}
func TestServiceDiscoveryMessage(t *testing.T) {
	{
		body := ``
		expected := `{"status":"success","code":"success","message":"","payload":{"pong":"pong"}}`
		resp, _ := testFunction("/sd_health", body)

		assert.Equal(t, resp, expected, "Cannot ping to the service discovery health check function.")
	}
}

func TestDatabase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Database didn't panic when the host is incorrect.")
		}
	}()

	resetDB := false

	os.Setenv("KITSVC_DATABASE_HOST", "xxxxxx")
	createDatabase(&resetDB)
}
