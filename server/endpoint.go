package server

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/context"
	"reflect"
)

type Endpoint interface {
	Name() string
	Handler() HandlerFunc
}

type DefaultEndpoint struct {
	name    string
	handler HandlerFunc
	req     interface{}
}

func (de *DefaultEndpoint) Name() string {
	return de.name
}

func (de *DefaultEndpoint) Handler() HandlerFunc {

	fn := func(ctx context.Context, req Request) (interface{}, error) {
		reqType := de.req
		err := proto.Unmarshal(req.Data(), reqType.(proto.Message))
		if err != nil {
			return nil, err
		}

		req.SetRequest(reqType)

		return de.handler(ctx, req)

	}

	return fn
}

func (de *DefaultEndpoint) Request() interface{} {
	return de.req
}

type ReflectEndpoint struct {
	name     string
	vhandler reflect.Value
	req      reflect.Type
}

func (de *ReflectEndpoint) Name() string {
	return de.name
}

func (de *ReflectEndpoint) Handler() HandlerFunc {
	return func(ctx context.Context, req Request) (interface{}, error) {

		b := req.Data()
		request := reflect.New((*de).req.Elem())

		protoData := request.Interface().(proto.Message)

		err := proto.Unmarshal(b, protoData)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshal data")

		}

		result := de.vhandler.Call([]reflect.Value{reflect.ValueOf(ctx), request})

		rsp := result[0].Interface()
		ok := true
		if !result[1].IsNil() {
			err, ok = result[1].Interface().(error)
			if !ok {
				return nil, fmt.Errorf("Non error return type from handler in second return argument")
			}
		}
		return rsp, err

	}
}
