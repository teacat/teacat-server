// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by a permissive BSD 3 Clause License
// that can be found in the license file.

package goes

import (
	"fmt"
	"net/http"
	"strconv"
)

// StreamWriter provides methods for writing events and metadata to an
// event stream.
type StreamWriter struct {
	client     *Client
	streamName string
}

// Append writes an event to the head of the stream.
//
// If the stream does not exist, it will be created.
//
// There are some special version numbers that can be provided.
// http://docs.geteventstore.com/http-api/3.7.0/writing-to-a-stream/
//
// -2 : The write should never conflict with anything and should always succeed.
//
// -1 : The stream should not exist at the time of writing. This write will create it.
//
// 0 : The stream should exist but it should be empty.
func (s *StreamWriter) Append(expectedVersion *int, events ...*Event) error {
	u := fmt.Sprintf("/streams/%s", s.streamName)
	req, err := s.client.NewRequest(http.MethodPost, u, events)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/vnd.eventstore.events+json")
	if expectedVersion != nil {
		req.Header.Set("ES-ExpectedVersion", strconv.Itoa(*expectedVersion))
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		if e, ok := err.(*ErrBadRequest); ok {
			return &ErrConcurrencyViolation{ErrorResponse: e.ErrorResponse}
		}
		return err
	}

	return nil
}

// WriteMetaData writes the metadata for a stream.
//
// The operation will replace the current stream metadata.
//
// For more information on stream metadata see:
// http://docs.geteventstore.com/http-api/3.7.0/stream-metadata/
//
// If the metadata was written successfully the error returned will be nil.
//
// If an error occurs the error returned may be an ErrUnauthorized, a
// ErrTemporarilyUnavailable or an ErrUnexpected if the error occurred during a
// http request to the server. In these cases, the *ErrorResponse will be available
// for inspection as an ErrorResponse field on the error.
// If an error occurred outside of the http request another type of error will be returned
// such as a *url.Error in cases where the streamwriter is unable to connect to the server.
func (s *StreamWriter) WriteMetaData(stream string, metadata interface{}) error {
	m := NewEvent("", "MetaData", metadata, nil)
	mURL, _, err := s.client.GetMetadataURL(stream)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(http.MethodPost, mURL, m)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/vnd.eventstore.events+json")

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
