package main

import "fmt"

type testE struct {
	msg string
}

func (svc service) CatchEvent(body map[string]interface{}, meta map[string]string) {
	fmt.Println(body["test"])
}
