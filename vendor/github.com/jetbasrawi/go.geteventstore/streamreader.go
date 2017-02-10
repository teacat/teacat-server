// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by a permissive BSD 3 Clause License
// that can be found in the license file.

package goes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jetbasrawi/go.geteventstore/atom"
)

// StreamReader provides methods for reading events and event metadata.
type StreamReader struct {
	streamName    string
	client        *Client
	version       int
	nextVersion   int
	index         int
	currentURL    string
	pageSize      int
	eventResponse *EventResponse
	feedPage      *atom.Feed
	lasterr       error
	loadFeedPage  bool
}

// Err returns any error that is raised as a result of a call to Next().
func (s *StreamReader) Err() error {
	return s.lasterr
}

// Version returns the current stream version of the reader.
func (s *StreamReader) Version() int {
	return s.version
}

// NextVersion is the version of the stream that will be returned by a call to Next().
func (s *StreamReader) NextVersion(version int) {
	s.nextVersion = version
}

// EventResponse returns the container for the event that is returned from a call to Next().
func (s *StreamReader) EventResponse() *EventResponse {
	return s.eventResponse
}

// Next gets the next event on the stream.
//
// Next should be treated more like a cursor over the stream rather than an
// enumerator over a collection of results. Individual events are retrieved
// on each call to Next().
//
// The boolean returned is intended to provide a convenient mechanism to
// to enumerate and process events, it should not be considered an indication
// of the status of a call to Next(). To undertand the outcomes of operations the
// stream's Err() field should be inspected. It is left to the user to determine
// under what conditions to exit the loop.
//
// When next is called, it will go to the eventstore and get a single event at the
// current reader's stream version.
func (s *StreamReader) Next() bool {
	s.lasterr = nil

	numEntries := 0
	if s.feedPage != nil {
		numEntries = len(s.feedPage.Entry)
	}

	// The feed page will be nil when the stream reader is first created.
	// The initial feed page url will be constructed based on the current
	// version number.
	if s.feedPage == nil {
		s.index = -1
		url, err := s.client.GetFeedPath(s.streamName, "forward", s.nextVersion, s.pageSize)
		if err != nil {
			s.lasterr = err
			return false
		}
		s.currentURL = url
	}

	// If the index is less than 0 load the previous feed page.
	// GetEventStore uses previous to point to more recent feed pages and uses
	// next to point to older feed pages. A stream starts at the most recent
	// event and ends at the oldest event.
	if s.index < 0 {
		if s.feedPage != nil {
			// Get the url for the previous feed page. If the reader is at the head
			// of the stream, the previous link in the feedpage will be nil.
			if l := s.feedPage.GetLink("previous"); l != nil {
				s.currentURL = l.Href
			}
		}

		//Read the feedpage at the current url
		f, _, err := s.client.ReadFeed(s.currentURL)
		if err != nil {
			s.lasterr = err
			return true
		}

		s.feedPage = f
		numEntries = len(f.Entry)
		s.index = numEntries - 1
	}

	//If there are no events returned at the url return an error
	if numEntries <= 0 {
		s.eventResponse = nil
		s.lasterr = &ErrNoMoreEvents{}
		return true
	}

	//There are events returned, get the event for the current version
	entry := s.feedPage.Entry[s.index]
	url := strings.TrimRight(entry.Link[1].Href, "/")
	e, _, err := s.client.GetEvent(url)
	if err != nil {
		s.lasterr = err
		return true
	}
	s.eventResponse = e
	s.version = s.nextVersion
	s.nextVersion++
	s.index--

	return true
}

// Scan deserializes event and event metadata into the types passed in
// as arguments e and m.
func (s *StreamReader) Scan(e interface{}, m interface{}) error {

	if s.lasterr != nil {
		return s.lasterr
	}

	if s.eventResponse == nil {
		return &ErrNoMoreEvents{}
	}

	if e != nil {
		data, ok := s.eventResponse.Event.Data.(*json.RawMessage)
		if !ok {
			return fmt.Errorf("Could not unmarshal the event. Event data is not of type *json.RawMessage")
		}

		if err := json.Unmarshal(*data, e); err != nil {
			return err
		}
	}

	if m != nil && s.EventResponse().Event.MetaData != nil {
		meta, ok := s.eventResponse.Event.MetaData.(*json.RawMessage)
		if !ok {
			return fmt.Errorf("Could not unmarshal the event. Event data is not of type *json.RawMessage")
		}

		if err := json.Unmarshal(*meta, &m); err != nil {
			return err
		}
	}

	return nil
}

// LongPoll causes the server to wait up to the number of seconds specified
// for results to become available at the URL requested.
//
// LongPoll is useful when polling the end of the stream as it returns quickly
// after new events are found which means that it is more responsive than polling
// over some arbitrary time interval. It also reduces the number of
// unproductive calls to the server polling for events when there are no new
// events to return.
//
// Setting the argument seconds to any integer value above 0 will cause the
// request to be made with ES-LongPoll set to that value. Any value 0 or below
// will cause the request to be made without ES-LongPoll and the server will not
// wait to return.
func (s *StreamReader) LongPoll(seconds int) {
	if seconds > 0 {
		s.client.SetHeader("ES-LongPoll", strconv.Itoa(seconds))
	} else {
		s.client.DeleteHeader("ES-LongPoll")
	}
}

// MetaData gets the metadata for a stream.
//
// Stream metadata is retured as an EventResponse.
//
// For more information on stream metadata see:
// http://docs.geteventstore.com/http-api/3.7.0/stream-metadata/
func (s *StreamReader) MetaData() (*EventResponse, error) {
	url, _, err := s.client.GetMetadataURL(s.streamName)
	if err != nil {
		return nil, err
	}
	ev, _, err := s.client.GetEvent(url)
	if err != nil {
		return nil, err
	}
	return ev, nil
}
