package errno

import (
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Err struct {
	Code       string
	Message    string
	StatusCode int
	Path       string
	Line       int
}

func (e *Err) Error() string {
	return e.Message
}

var (
	errs = map[string]*Err{
		"Bind": {Code: "BIND_ERR", Message: "Error occurred while binding the request body to the struct.", StatusCode: http.StatusBadRequest},
	}
	//ErrBind = &Err{Code: "BIND_ERR", Message: "Error occurred while binding the request body to the struct.", StatusCode: http.StatusBadRequest}
)

func Error(err string) *Err {
	if val, ok := errs[err]; ok {
		_, fn, line, _ := runtime.Caller(1)

		val.Path = strings.Replace(fn, os.Getenv("GOPATH"), "", -1)
		val.Line = line

		return val
	}
	return nil
}
