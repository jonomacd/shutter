package shutter

import (
	"github.com/gogo/protobuf/proto"
	message "github.com/jonomacd/shutter/proto"
	"golang.org/x/net/context"
)

type Request interface {
	context.Context

	GetHeader(string) string
	SetHeader(string, string)

	Data() []byte
	Request() interface{}
}

type DefaultRequest struct {
	context.Context

	headers  map[string]string
	service  string
	endpoint string
	body     []byte
	req      interface{}
}

func NewRequest(service, endpoint string, body []byte, reqType interface{}, headers []message.Keyvalue) Request {

	req := &DefaultRequest{
		service:  service,
		endpoint: endpoint,
		body:     body,
	}

	// Testing... need to deal with errors
	proto.Unmarshal(body, reqType.(proto.Message))
	req.req = reqType

	for _, kv := range headers {
		req.SetHeader(kv.Key, kv.Value)
	}

	return req

}

func (dr *DefaultRequest) GetHeader(key string) string {
	return dr.headers[key]
}

func (dr *DefaultRequest) SetHeader(key, value string) {
	dr.headers[key] = value
}

func (dr *DefaultRequest) Data() []byte {
	return dr.body
}

func (dr *DefaultRequest) Request() interface{} {
	return dr.req
}
