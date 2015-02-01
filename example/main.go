package main

import (
	"fmt"
	"github.com/jonomacd/shutter"
	"github.com/jonomacd/shutter/client"

	message "github.com/jonomacd/shutter/proto"
)

var (
	serviceName string = "uniqueServiceName"
	endpoint    string = "firstEndpoint"
)

func main() {
	shutter.InitializeService(serviceName)

	shutter.ServiceRegistry.Register(endpoint, TestHandler, &message.Keyvalue{})

	// Run Client
	send()

}

func send() {

	toSend := &message.Keyvalue{
		Key:   "foo",
		Value: "bar",
	}

	toGet := &message.Keyvalue{}

	err := client.Request(serviceName, endpoint, toSend, toGet)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Recieved: %s\n", toGet.Value)

}

func TestHandler(req shutter.Request) (interface{}, error) {
	r := req.Request().(*message.Keyvalue)

	r.Value = r.Value + " From Handler " + r.Key

	return r, nil
}
