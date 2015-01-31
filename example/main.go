package main

import (
	"fmt"
	"github.com/jonomacd/shutter"
	"time"

	"github.com/gogo/protobuf/proto"
	message "github.com/jonomacd/shutter/proto"
	"gopkg.in/project-iris/iris-go.v1"
)

func main() {
	shutter.InitializeService("test")

	shutter.ServiceRegistry.Register("first", func(req shutter.Request) ([]byte, error) {
		fmt.Printf("Got It %s\n", string(req.Data()))
		return []byte("{}"), nil
	})

	shutter.ServiceRegistry.Register("first1", func(req shutter.Request) ([]byte, error) {
		fmt.Printf("Got It1 %s\n", string(req.Data()))
		return []byte("{}1"), nil
	})
	fmt.Println("sending")
	send()

	select {}
}

func send() {
	req := &message.Request{}
	req.Body = []byte("Hi There")
	req.Endpoint = "first"
	req.Originator = "me"

	b, err := proto.Marshal(req)
	fmt.Println(err)

	conn, err := iris.Connect(55555)
	if err != nil {
		panic(err)
	}

	cruft, err := conn.Request("test", b, time.Second*10)
	if err != nil {
		panic(err)
	}
	fmt.Println("got back", string(cruft))

	req = &message.Request{}
	req.Body = []byte("Hi There1")
	req.Endpoint = "first1"
	req.Originator = "me1"

	b, err = proto.Marshal(req)
	fmt.Println(err)

	cruft, err = conn.Request("test", b, time.Second*10)
	if err != nil {
		panic(err)
	}
	fmt.Println("got back1", string(cruft))
}
