package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/TeaMeow/KitSvc/service"
	"github.com/go-kit/kit/endpoint"
)

func makeUppercaseEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.UppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return service.UppercaseResponse{v, err.Error()}, nil
		}
		return service.UppercaseResponse{v, ""}, nil
	}
}

func makeLowercaseEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.LowercaseRequest)
		v, err := svc.Lowercase(req.S)
		if err != nil {
			return service.LowercaseResponse{v, err.Error()}, nil
		}
		return service.LowercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.CountRequest)
		v := svc.Count(req.S)
		return service.CountResponse{v}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLowercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.LowercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.CountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
