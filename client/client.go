package client

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	message "github.com/jonomacd/shutter/proto"
	"gopkg.in/project-iris/iris-go.v1"
)

var (
	Port int = 55555
	Conn *iris.Connection
)

func init() {
	var err error
	Conn, err = iris.Connect(Port)
	if err != nil {
		// TERRIBLE, need better err handling
		panic(err)
	}

}

func Request(service, endpoint string, req proto.Message, rsp proto.Message) error {
	wireReq := &message.Request{
		Endpoint:   endpoint,
		Originator: "",
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	wireReq.Body = b

	// Todo error checking
	wireb, _ := proto.Marshal(wireReq)

	// Todo change timeout
	rspb, err := Conn.Request(service, wireb, time.Second*10)
	if err != nil {
		panic(err)
	}

	wireRsp := &message.Response{}

	proto.Unmarshal(rspb, wireRsp)

	if wireRsp.Type == "error" {
		// Todo, deal with errors better.  Build specific types based on the code
		return fmt.Errorf("%s :: %s", wireRsp.Err.ErrorText, wireRsp.Err.Code)
	}

	err = proto.Unmarshal(wireRsp.Body, rsp)
	return err
}
