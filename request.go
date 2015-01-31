package shutter

import (
	message "github.com/jonomacd/shutter/proto"
	"golang.org/x/net/context"
)

type Request interface {
	context.Context

	GetHeader(string) string
	SetHeader(string, string)

	Data() []byte
}

type DefaultRequest struct {
	context.Context

	headers  map[string]string
	service  string
	endpoint string
	body     []byte
}

func NewRequest(service, endpoint string, body []byte, headers []message.Keyvalue) Request {

	req := &DefaultRequest{
		service:  service,
		endpoint: endpoint,
		body:     body,
	}

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

func (dr *DefaultRequest) Send(service, endpoint string, data []byte) {}
