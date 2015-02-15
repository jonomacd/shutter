# Shutter
Microservice library on top of Iris

# Example With Reflection

```
package main

import (
	"log"

	"github.com/jonomacd/shutter/client"
	"github.com/jonomacd/shutter/server"
	"golang.org/x/net/context"

	message "github.com/jonomacd/shutter/proto"
)

var (
	serviceName string = "uniqueServiceName"
	endpoint    string = "hello"
)

func main() {
	// Connect the server, Mostly this is so you can name your service
	server.InitializeService(serviceName)

	// Register an endpoint, the endpoint name will be a lower case version of the fuction name.
	server.SimpleRegister(Hello)

	// Run Client
	// Construct response object
	toGet := &message.Keyvalue{}

	// Send the request
	client.Request(serviceName, endpoint, &message.Keyvalue{Key: "foo", Value: "bar"}, toGet)

	log.Printf("Received: %s\n", toGet)

}

func Hello(ctx context.Context, m *message.Keyvalue) (*message.Keyvalue, error) {
	return &message.Keyvalue{
		Key:   "hello",
		Value: m.Value + " From hello Handler " + m.Key,
	}, nil
}
```

# Example Without Reflection

```
package main

import (
	"fmt"
	"log"

	"github.com/jonomacd/shutter/client"
	"github.com/jonomacd/shutter/server"
	"golang.org/x/net/context"

	message "github.com/jonomacd/shutter/proto"
)

var (
	serviceName string = "uniqueServiceName"
	endpoint    string = "firstEndpoint"
)

func main() {
	// Connect the server
	server.InitializeService(serviceName)

	// Register an endpoint
	server.Register(endpoint, TestHandler, &message.Keyvalue{})

	// Run Client
	send()

}

func send() {
	// Initialize a global client (can mint our own client but this is mostly for convience)
	client.InitializeClient()

	// Construct our request and response objects (happen to be the same type because I am reusing a proto definition)
	toSend := &message.Keyvalue{
		Key:   "foo",
		Value: "bar",
	}
	toGet := &message.Keyvalue{}

	// Send the request
	err := client.Request(serviceName, endpoint, toSend, toGet)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Printf("Received: %s\n", toGet)

}

// TestHandler is a handler function for our test server
func TestHandler(ctx context.Context, req server.Request) (interface{}, error) {
	r, ok := req.Request().(*message.Keyvalue)
	if !ok {
		return nil, fmt.Errorf("Bad Type")
	}

	response := &message.Keyvalue{
		Key:   "something",
		Value: r.Value + " From Handler " + r.Key,
	}

	return response, nil
}
```
