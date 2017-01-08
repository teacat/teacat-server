package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUppercase(t *testing.T) {
	{
		body := `{"s": "test"}`
		expected := `{"status":"success","code":"success","message":"","payload":{"v":"TEST"}}`
		resp, _ := testFunction("POST", "/uppercase", body)

		assert.Equal(t, expected, resp, "Uppercase() cannot convert the string to uppercase.")
	}
	{
		body := `{"s":`
		expected := `{"status":"error","code":"error","message":"Cannot parse the JSON content.","payload":null}`
		resp, _ := testFunction("POST", "/uppercase", body)

		assert.Equal(t, expected, resp, "Uppercase() cannot tell when the parse error occurred.")
	}
	{
		body := `{"s": ""}`
		expected := `{"status":"error","code":"string_empty","message":"The string is empty.","payload":null}`
		resp, _ := testFunction("POST", "/uppercase", body)

		assert.Equal(t, expected, resp, "Uppercase() cannot tell the error.")
	}
}

func TestCount(t *testing.T) {
	{
		body := `{"s": "test"}`
		expected := `{"status":"success","code":"success","message":"","payload":{"v":4}}`
		resp, _ := testFunction("POST", "/count", body)

		assert.Equal(t, expected, resp, "Count() cannot count the length of the string.")
	}
	{
		body := `{"s":`
		expected := `{"status":"error","code":"error","message":"Cannot parse the JSON content.","payload":null}`
		resp, _ := testFunction("POST", "/count", body)

		assert.Equal(t, expected, resp, "Count() cannot tell when the parse error occurred.")
	}
}
