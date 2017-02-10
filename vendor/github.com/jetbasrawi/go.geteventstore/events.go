// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by a permissive BSD 3 Clause License
// that can be found in the license file.

package goes

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/jetbasrawi/go.geteventstore/internal/uuid"
)

// EventResponse encapsulates the response for an event reflecting the atom
// response returned from the server which contains data in addition to the
// actual event when requested as content type application/vnd.eventstore.atom+json
//
// For more information on the server response see:
// http://docs.geteventstore.com/http-api/latest/reading-streams/
type EventResponse struct {
	Title   string
	ID      string
	Updated TimeStr
	Summary string
	Event   *Event
}

// PrettyPrint renders an indented json view of the EventResponse.
func (e *EventResponse) PrettyPrint() string {

	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)

}

// EventAtomResponse is used internally to unmarshall the raw response
type EventAtomResponse struct {
	Title   string      `json:"title"`
	ID      string      `json:"id"`
	Updated TimeStr     `json:"updated"`
	Summary string      `json:"summary"`
	Content interface{} `json:"content"`
}

// PrettyPrint renders and indented json view of the eventAtomResponse
func (e *EventAtomResponse) PrettyPrint() string {

	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Event encapsulates the data of an eventstore event.
//
// EventStreamID is the id returned in the event atom response.
// EventNumber represents the stream version for this event.
// EventType describes the event type.
// EventID is the guid of the event.
// Data contains the data of the event.
// Links contains the urls of the event on the evenstore
// MetaData contains the metadata for the event.
type Event struct {
	EventStreamID string      `json:"eventStreamId,omitempty"`
	EventNumber   int         `json:"eventNumber,omitempty"`
	EventType     string      `json:"eventType,omitempty"`
	EventID       string      `json:"eventId,omitempty"`
	Data          interface{} `json:"data"`
	Links         []Link      `json:"links,omitempty"`
	MetaData      interface{} `json:"metadata,omitempty"`
}

// PrettyPrint renders an indented json view of the Event object.
func (e *Event) PrettyPrint() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Link encapsulates url data for events.
type Link struct {
	URI      string `json:"uri"`
	Relation string `json:"relation"`
}

// TimeStr is a type used to format feed dates.
type TimeStr string

// Time returns a TimeStr version of the time.Time argument t.
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

// NewEvent creates a new event object.
//
// If an empty eventId is provided a new uuid will be generated automatically
// and retured in the event.
// If an empty eventType is provided the eventType will be set to the
// name of the type provided.
// data and meta can be nil.
func NewEvent(eventID, eventType string, data interface{}, meta interface{}) *Event {
	e := &Event{}

	e.EventID = eventID
	if eventID == "" {
		e.EventID = NewUUID()
	}

	e.EventType = eventType
	if eventType == "" {
		e.EventType = typeOf(data)
	}

	e.Data = data
	e.MetaData = meta
	return e
}

// NewUUID returns a new V4 uuid as a string.
func NewUUID() string {
	return uuid.NewV4().String()
}

// typeOf is a helper to get the names of types.
func typeOf(i interface{}) string {
	if i == nil {
		return ""
	}
	return reflect.TypeOf(i).Elem().Name()
}
