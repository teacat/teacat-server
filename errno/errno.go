package errno

import (
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
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

func Error(code string) *Err {
	if val, ok := errs[code]; ok {
		_, fn, line, _ := runtime.Caller(1)

		val.Code = code
		val.Path = strings.Replace(fn, os.Getenv("GOPATH"), "", -1)
		val.Line = line

		return val
	}
	return nil
}

func Abort(code string, err error, c *gin.Context) {
	c.Error(err)
	c.Error(Error(code))
	c.Abort()
}
