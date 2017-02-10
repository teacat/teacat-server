// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by a permissive BSD 3 Clause License
// that can be found in the license file.

package goes

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jetbasrawi/go.geteventstore/atom"
)

// Response encapsulates HTTP responses from the server.
//
// A Response object contains the raw http response, the status code returned
// and the status message returned.
//
// A Response object is returned from all methods on the client that interact
// with the the eventstore.
// It is intended to provide access to data about the response in case
// the user wants to inspect the response such as the status or the raw http
// response.
type Response struct {
	*http.Response
	Status     string
	StatusCode int
}

// ErrorResponse encapsulates data about an interaction with the eventstore that
// produced an HTTP error.
//
// An ErrorResponse embeds the raw *http.Response and provides access to the raw
// http.Request that resulted in an error.
// Status contains the status message returned from the server.
// StatusCode contains the status code returned from the server.
type ErrorResponse struct {
	*http.Response
	Request    *http.Request
	Status     string
	StatusCode int
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Status)
}

type basicAuthCredentials struct {
	Username string
	Password string
}

// Client is the interface that the client should implement
// type Client interface {
// 	NewStreamReader(streamName string) *StreamReader
// 	NewStreamWriter(streamName string) *StreamWriter
// 	SetBasicAuth(username, password string)
// 	NewRequest(method, urlString string, body interface{}) (*http.Request, error)
// 	Do(req *http.Request, v io.Writer) (*Response, error)
// 	GetEvent(url string) (*EventResponse, *Response, error)
// 	GetMetadataURL(stream string) (string, *Response, error)
// 	ReadFeed(url string) (*atom.Feed, *Response, error)
// 	SetHeader(key, value string)
// 	DeleteHeader(key string)
// }

// Client is a handle for an eventstore server.
//
// The client is used to store connection details such as server URL,
// basic authentication credentials and headers.
//
// The client also provides methods interacting with the eventstore.
//
// In general, the StreamReader and StreamWriter should be used to
// interact with the eventstore. These methods further abstract methods
// on the client, however you can also directly use methods on the client
// to interact with the eventstore if you want to create some custom behaviour.
type Client struct {
	client      *http.Client
	baseURL     *url.URL
	credentials *basicAuthCredentials
	headers     map[string]string
}

// NewClient returns a new client.
//
// httpClient will usually be nil and the client will use the http.DefaultClient.
// Should you want to implement behaviours at the transport level you can provide
// your own *http.Client
//
// serverURL is the full URL to your eventstore server including protocol scheme and
// port number.
func NewClient(httpClient *http.Client, serverURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:  httpClient,
		baseURL: baseURL,
		headers: make(map[string]string),
	}
	return c, nil
}

// NewStreamReader returns a new *StreamReader.
func (c *Client) NewStreamReader(streamName string) *StreamReader {
	return &StreamReader{
		streamName: streamName,
		client:     c,
		version:    -1,
		pageSize:   20,
	}
}

// NewStreamWriter returns a new *StreamWriter.
func (c *Client) NewStreamWriter(streamName string) *StreamWriter {
	return &StreamWriter{
		client:     c,
		streamName: streamName,
	}
}

// SetBasicAuth sets the credentials for requests.
//
// Credentials will be read from the client before each request.
func (c *Client) SetBasicAuth(username, password string) {
	c.credentials = &basicAuthCredentials{
		Username: username,
		Password: password,
	}
}

// GetEvent reads a single event from the eventstore.
//
// The event response will be nil in an error case.
// *Response may be nil if an error occurs before the http request. Otherwise
// it will contain the raw http response and status.
// If an error occurs during the http request an *ErrorResponse will be returned
// as the error. The *ErrorResponse will contain the raw http response and status
// and a description of the error.
func (c *Client) GetEvent(url string) (*EventResponse, *Response, error) {

	r, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	r.Header.Set("Accept", "application/vnd.eventstore.atom+json")

	var b bytes.Buffer
	resp, err := c.Do(r, &b)
	if err != nil {
		return nil, resp, err
	}

	if b.String() == "{}" {
		return nil, resp, nil
	}
	var raw json.RawMessage
	er := &EventAtomResponse{Content: &raw}
	err = json.NewDecoder(bytes.NewReader(b.Bytes())).Decode(er)
	if err == io.EOF {
		return nil, resp, nil
	}

	if err != nil {
		return nil, resp, err
	}

	var d json.RawMessage
	var m json.RawMessage
	ev := &Event{Data: &d, MetaData: &m}

	err = json.Unmarshal(raw, ev)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, resp, err
	}

	e := EventResponse{}
	e.Title = er.Title
	e.ID = er.ID
	e.Updated = er.Updated
	e.Summary = er.Summary
	e.Event = ev

	return &e, resp, nil
}

// ReadFeed reads the atom feed for a stream and returns an *atom.Feed.
//
// The feed object returned may be nil in case of an error.
// The *Response may also be nil if the error occurred before the http request.
// If the error occurred after the http request, the *Response will contain the
// raw http response and status.
// If the error occurred during the http request an *ErrorResponse will be returned
// and this will also contain the raw http request and status and an error message.
func (c *Client) ReadFeed(url string) (*atom.Feed, *Response, error) {

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/atom+xml")

	var b bytes.Buffer
	resp, err := c.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}
	feed := &atom.Feed{}
	err = xml.NewDecoder(bytes.NewReader(b.Bytes())).Decode(feed)
	if err != nil {
		return nil, resp, err
	}

	return feed, resp, nil
}

// GetFeedPath returns the path for a feedpage
//
// Valid directions are "forward" and "backward".
//
// To get the path to the head of the stream, pass a negative integer in the version
// argument and "backward" as the direction.
func (c *Client) GetFeedPath(stream, direction string, version int, pageSize int) (string, error) {
	ps := pageSize

	dir := ""
	switch direction {
	case "forward", "backward":
		dir = direction
	default:
		return "", fmt.Errorf("Invalid Direction %s. Allowed values are \"forward\" or \"backward\" \n", direction)
	}

	v := "head"
	if version >= 0 {
		v = strconv.Itoa(version)
	}

	if v == "head" && dir == "forward" {
		return "", fmt.Errorf("Invalid Direction (%s) and version (head) combination.\n", direction)
	}

	return fmt.Sprintf("/streams/%s/%s/%s/%d", stream, v, dir, ps), nil
}

// GetMetadataURL gets the url for the stream metadata.
// according to the documentation the metadata url should be acquired through
// a query to the stream feed as the authors of GetEventStore reserve the right
// to change the url.
// http://docs.geteventstore.com/http-api/latest/stream-metadata/
func (c *Client) GetMetadataURL(stream string) (string, *Response, error) {

	url, err := c.GetFeedPath(stream, "forward", 0, 1)
	if err != nil {
		return "", nil, err
	}

	f, resp, err := c.ReadFeed(url)
	if err != nil {
		return "", resp, err
	}
	for _, v := range f.Link {
		if v.Rel == "metadata" {
			return v.Href, resp, nil
		}
	}
	return "", resp, nil
}

// SetHeader adds a header to the collection of headers that will be used on http requests.
//
// Any headers that are set on the client will be included in requests to the eventstore.
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// DeleteHeader deletes a header from the collection of headers.
func (c *Client) DeleteHeader(key string) {
	delete(c.headers, key)
}

// DeleteStream will delete a stream
//
// Streams may be soft deleted or hard deleted.
//
// Soft deleting a stream means that you can later recreate it simply by appending
// events to it.
//
// Hard deleting a stream means that it has permanently been deleted and can never
// be recreated.
//
// http://docs.geteventstore.com/http-api/3.8.0/deleting-a-stream/
func (c *Client) DeleteStream(streamName string, hardDelete bool) (*Response, error) {

	url := fmt.Sprintf("/streams/%s", streamName)

	req, err := c.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	if hardDelete {
		req.Header.Set("ES-HardDelete", "true")
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil

}

// NewRequest creates a new *http.Request that can be used to execute requests to the
// server using the client.
func (c *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {

	url, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if !url.IsAbs() {
		url = c.baseURL.ResolveReference(url)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.credentials != nil {
		req.SetBasicAuth(c.credentials.Username, c.credentials.Password)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Do executes requests to the server.
//
// The response body is copied into v if v is not nil.
// Returns a *Response that wraps the http.Response returned from the server.
// The response body is available in the *Response in case the consumer wishes
// to process it in some way rather than read if from the argument v
func (c *Client) Do(req *http.Request, v io.Writer) (*Response, error) {

	// keep is a copy of the request body that will be returned
	// with the response for diagnostic purposes.
	// send will be used to make the request.
	var keep, send io.ReadCloser

	if req.Body != nil {
		if buf, err := ioutil.ReadAll(req.Body); err == nil {
			keep = ioutil.NopCloser(bytes.NewReader(buf))
			send = ioutil.NopCloser(bytes.NewReader(buf))
			req.Body = send
		}
	}

	// An error is returned if caused by client policy (such as CheckRedirect),
	// or if there was an HTTP protocol error. A non-2xx response doesn't cause
	// an error.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Create a *Response to wrap the http.Response
	response := newResponse(resp)

	// After the request has been made the req.Body will be unreadable.
	// assign keep to the request body so that it can be returned in the
	// response for diagnostic purposes.
	if keep != nil {
		req.Body = keep
	}

	// If the request returned an error status checkResponse will return an
	// *errorResponse containing the original request, status code and status message
	err = getError(resp, req)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	// When handling post requests v will be nil
	if v != nil {
		io.Copy(v, resp.Body)
	}

	return response, nil
}

// getError inspects the HTTP response and constructs an appropriate error if
// the response was an error.
func getError(r *http.Response, req *http.Request) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	errorResponse.Status = r.Status
	errorResponse.StatusCode = r.StatusCode
	errorResponse.Request = req

	switch r.StatusCode {
	case http.StatusBadRequest:
		return &ErrBadRequest{ErrorResponse: errorResponse}
	case http.StatusUnauthorized:
		return &ErrUnauthorized{ErrorResponse: errorResponse}
	case http.StatusServiceUnavailable:
		return &ErrTemporarilyUnavailable{ErrorResponse: errorResponse}
	case http.StatusNotFound:
		return &ErrNotFound{ErrorResponse: errorResponse}
	case http.StatusGone:
		return &ErrDeleted{ErrorResponse: errorResponse}
	default:
		return &ErrUnexpected{ErrorResponse: errorResponse}
	}
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.Status = r.Status
	response.StatusCode = r.StatusCode
	return response
}
