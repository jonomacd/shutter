package server

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	message "github.com/jonomacd/shutter/proto"
)

type ServerError struct {
	Code      string
	ErrorText string
}

func (se *ServerError) Error() string {
	return fmt.Sprintf("Code: %s, Error: %s", se.Code, se.Error)
}

func WireError(err error) []byte {

	se, ok := err.(*ServerError)
	if !ok {
		se = &ServerError{
			Code:      "unknown",
			ErrorText: err.Error(),
		}
	}

	protoRsp := &message.Response{
		Type: "error",
		Err: message.Error{
			Code:      se.Code,
			ErrorText: se.ErrorText,
		},
	}
	b, _ := proto.Marshal(protoRsp)

	return b

}
