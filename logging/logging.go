package logging

import (
	"github.com/TeaMeow/KitSvc/service"
	"github.com/go-kit/kit/log"
)

// logmw is the middleware of the logger.
type Middleware struct {
	Logger  log.Logger
	Service service.Service
}

// loggingMiddleware recevie a logger and returns a new logging middleware.
func CreateMiddleware(logger log.Logger) service.Middleware {
	return func(next service.Service) service.Service {
		return Middleware{Logger: logger, Service: next}
	}
}
