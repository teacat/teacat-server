package errno

import (
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

// Err represents an error, `Path`, `Line`, `Code` will be automatically filled.
type Err struct {
	Code       string
	Message    string
	StatusCode int
	Path       string
	Line       int
}

// Error returns the error message.
func (e *Err) Error() string {
	return e.Message
}

// Error returns the error struct by the error code.
func Error(code string) *Err {
	if val, ok := errs[code]; ok {
		_, fn, line, _ := runtime.Caller(1)

		// Fill the error occurred path, line, code.
		val.Code = code
		val.Path = strings.Replace(fn, os.Getenv("GOPATH"), "", -1)
		val.Line = line

		return val
	}
	return nil
}

// Abort the current request with the specified error code.
func Abort(code string, err error, c *gin.Context) {
	c.Error(err)
	c.Error(Error(code))
	c.Abort()
}
