// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by a permissive BSD 3 Clause License
// that can be found in the license file.

package goes

// ErrNoMoreEvents is returned when there are no events to return
// from a request to a stream.
type ErrNoMoreEvents struct{}

func (e ErrNoMoreEvents) Error() string {
	return "There are no more events to load."
}

// ErrNotFound is returned when a stream is not found.
type ErrNotFound struct {
	ErrorResponse *ErrorResponse
}

func (e ErrNotFound) Error() string {
	return "The stream does not exist."
}

// ErrDeleted is returned when a request is made to a stream that
// has been hard deleted.
type ErrDeleted struct {
	ErrorResponse *ErrorResponse
}

func (e ErrDeleted) Error() string {
	return "The stream has was deleted."
}

// ErrUnauthorized is returned when a request to the eventstore is
// not authorized
type ErrUnauthorized struct {
	ErrorResponse *ErrorResponse
}

func (e ErrUnauthorized) Error() string {
	return "You are not authorised to access the stream or the stream does not exist."
}

// ErrTemporarilyUnavailable is returned when the server returns ServiceUnavailable.
//
// This error may be returned if a request is made to the server during startup. When
// the server starts up initially and the client is completely unable to connect to the
// server a *url.Error will be returned. Once the server is up but not ready to serve
// requests a ServiceUnavailable error will be returned for a brief period.
type ErrTemporarilyUnavailable struct {
	ErrorResponse *ErrorResponse
}

func (e ErrTemporarilyUnavailable) Error() string {
	return "Server Is Not Ready"
}

// ErrUnexpected is returned when a request to the eventstore returns an error that
// is not explicitly represented by a goes Error type such as UnauthorisedError or
// ErrNotFound
type ErrUnexpected struct {
	ErrorResponse *ErrorResponse
}

func (e ErrUnexpected) Error() string {
	return "An unexpected error occurred."
}

// ErrBadRequest is returned when the server returns a bad request error
type ErrBadRequest struct {
	ErrorResponse *ErrorResponse
}

func (e ErrBadRequest) Error() string {
	return "Bad request."
}

// ErrConcurrencyViolation is returned when the expected version does not match
// the stream version when writing to an event stream.
type ErrConcurrencyViolation struct {
	ErrorResponse *ErrorResponse
}

func (e ErrConcurrencyViolation) Error() string {
	return "Concurrency Error."
}
