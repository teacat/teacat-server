package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V string `json:"v"`
}

func makeUppercaseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v}, err
		}
		return uppercaseResponse{v}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// Uppercase logs the informations about the Uppercase function of the service.
func (mw LoggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

// Uppercase records the instrument about the Uppercase function of the service.
func (mw InstrumentingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

// Uppercase converts the string to uppercase.
func (svc service) Uppercase(s string) (string, error) {

	res, err := svc.Model.ToUpper(s)
	if err != nil {
		return "", err
	}

	return res, nil
}
