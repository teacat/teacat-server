#Go.GetEventStore [![license](https://img.shields.io/badge/license-BSD-blue.svg?maxAge=2592000)](https://github.com/jetbasrawi/go.geteventstore/blob/master/LICENSE.md) [![Go Report Card](https://goreportcard.com/badge/github.com/jetbasrawi/go.geteventstore)](https://goreportcard.com/report/github.com/jetbasrawi/go.geteventstore) [![GoDoc](https://godoc.org/github.com/jetbasrawi/go.geteventstore?status.svg)](https://godoc.org/github.com/jetbasrawi/go.geteventstore)

##A Golang client for EventStore 3.x HTTP API. 
Go.GetEventStore is a http client for [GetEventStore](https://geteventstore.com) written in Go. The 
client abstracts interaction with the GetEventStore HTTP API providing easy to use features 
for reading and writing of events and event metadata.

##CQRS reference implementation
An full example CQRS implementation using go.geteventstore can be found at [go.cqrs](https://github.com/jetbasrawi/go.cqrs)

###Supported features
| Feature | Description |
|---------|-------------|
| **Write Events & Event Metadata** | Writing single and multiple events to a stream. Optionally expected version can be provided if you want to use optimistic concurrency features of the eventstore. |
| **Read Events & Event Metadata** | Reading events & event metadata from a stream. |
| **Read & Write Stream Metadata** | Read and writing stream metadata. |
| **Basic Authentication** | |
| **Long Poll** | Long Poll allows the client to listen at the head of a stream for new events. |
| **Soft & Hard Delete Stream** | |
| **Catch Up Subscription** | Using long poll with a StreamReader provides an effective catch up subscription. |
| **Serialization & Deserialization of Events** | The package handles serialization and deserialization of your application events to and from the eventstore. |
| **Reading Stream Atom Feed** | The package provides methods for reading stream Atom feed pages, returning a fully typed struct representation. |
| **Setting Optional Headers** | Optional headers can be added and removed. |

Below are some code examples giving a summary view of how the client works. To learn to use 
the client in more detail, heavily commented example code can be found in the examples directory.

###Get the package
```
    $ go get github.com/jetbasrawi/go.geteventstore
```

###Import the package
```go 
    import "github.com/jetbasrawi/go.geteventstore"
```

###Create a new client

```go
    client, err := goes.NewClient(nil, "http://youreventstore:2113")
	if err != nil {
		log.Fatal(err)
	}

```

###Set basic authentication

If required, you can set authentication on the client. Credentials can be changed at any time.
Requests are made with the credentials that were set last or none if none are set.

```go

    client.SetBasicAuth("admin", "changeit")

```

###Write events and event Metadata

Writing events and event metadata are supported via the StreamWriter. 

```go

    // Create your event
	myEvent := &FooEvent{
		FooField: "Lorem Ipsum",
		BarField: "Dolor Sit Amet",
		BazField: 42,
	}

    // Create your metadata type
    myEventMeta := make(map[string]string)
	myEventMeta["Foo"] = "consectetur adipiscing elit"

    // Wrap your event and event metadata in a goes.Event
	myGoesEvent := goes.NewEvent(goes.NewUUID(), "FooEvent", myEvent, myEventMeta)

    // Create a new StreamWriter
    writer := client.NewStreamWriter("FooStream")

    // Write the event to the stream, here we pass nil as the expectedVersion as we 
    // are not wanting to flag concurrency errors
    err := writer.Append(nil, myGoesEvent)
    if err != nil {
        // Handle errors
    }

```

###Read events

Reading events using the goes.StreamReader loosely follows the iterator idiom used in 
the "database/sql" package and other go libraries that deal with databases. This idiom 
works very well for reading events and provides an easy way to control the rate at which 
events are returned and to poll at the head of a stream.

An example of how to read all the events in a stream and then exit can be found in the 
read and write events example.

An example of how to read up to the head of the stream and then continue to listen for new 
events can be found in the longpoll example.

```go 

    // Create a new goes.StreamReader
    reader := client.NewStreamReader("FooStream")
    // Call Next to get the next event
    for reader.Next() {
        // Check if the call resulted in an error. 
        if reader.Err() != nil {
            // Handle errors
        }
        // If the call did not result in an error then an event was returned
        // Create the application types to hold the deserialized even data and meta data
        fooEvent := FooEvent{}
        fooMeta := make(map[string]string)
        // Call scan to deserialize the event data and meta data into your types
        err := reader.Scan(&fooEvent, &fooMeta)
        if err != nil {
            // Handle errors that occured during deserialization
        }
    }

```

###Long polling head of a stream

LongPoll provides an easy and efficient way to poll a stream listening for new events. 
The server will wait the specified amount of time or until new events are available on 
a stream. 

```go 

    reader := client.NewStreamReader("FooStream")
    for reader.Next() {
        if reader.Err() != nil {
            // When there are no more event in the stream, set LongPoll. 
            // The server will wait for 15 seconds in this case or until
            // events become available on the stream.
            if e, ok := reader.Err().(*goes.ErrNoMoreEvents); ok {
                reader.LongPoll(15)
            }
        } else {
            fooEvent := FooEvent{}
            _ := reader.Scan(&fooEvent, &fooMeta)
        }
    }

```

A more detailed example of using LongPoll can be found in the examples directory.

###Deleting streams

The client supports both soft delete and hard delete of event streams. 

```go

    // Soft delete or hard delete is specified by a boolean argument
    // here foostream will be soft deleted. 
    resp, err := client.DeleteStream("foostream", false)

```

Example code for deleting streams can be found in the examples directory.

###Direct use of the client

The StreamReader and StreamWriter types are the easiest way to read and write events. If you would like to implement
some other logic around reading events, some methods are available on the Client type. An example is included 
demonstrating how to use some of these methods.

###Setting optional headers 
Most of the optional headers are are included implicitly such as ES-LongPoll, ES-ExpectedVersion & ES-HardDelete when 
using LongPoll on the StreamReader or when appending events or deleting streams. However should you wish to use
any of the others the can be set explicitly on the client. The only other ones you might want to use are ES-ResolveLinkTo, 
ES-RequiresMaster or ES-TrustedAuth.

```go
    // Setting a header will mean all subsequent requests will include the header
    client.SetHeader("ES-ResolveLinkTo", "false")

    // Deleting a header means that it will not be included in any subsequent requests.
    client.DeleteHeader("ES-ResolveLinkTo")

```
###Running the Unit Tests
To keep the library lightweight and easy to use, I have tried not to have any dependencies on other 
packages. To use the package there are no dependencies, to run the unit tests however, the package does require 
some dependencies.

You will need the go.geteventstore.testfeed package.

```
    $ go get github.com/jetbasrawi/go.geteventstore.testfeed

```
You will also need the Check.V1 package which is used for assertions in unit tests.

```
    $ go get gopkg.in/check.v1
```

After this you should be able to run the tests as normal for any Golang project with unit tests.

``` 
    $ go test

```

###Feedback and requests welcome

This is a pretty new piece of work and criticism, comments or complements are most welcome.





