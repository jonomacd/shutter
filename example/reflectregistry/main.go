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
