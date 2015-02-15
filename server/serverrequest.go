package server

import (
	message "github.com/jonomacd/shutter/proto"
)

type Request interface {
	GetHeader(string) string
	SetHeader(string, string)

	Data() []byte
	Request() interface{}
	SetRequest(i interface{})
}

type DefaultRequest struct {
	headers  map[string]string
	service  string
	endpoint string
	body     []byte
	req      interface{}
}

func NewRequest(service, endpoint string, body []byte, headers []message.Keyvalue) (Request, error) {

	req := &DefaultRequest{
		service:  service,
		endpoint: endpoint,
		body:     body,
	}

	for _, kv := range headers {
		req.SetHeader(kv.Key, kv.Value)
	}

	return req, nil

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

func (dr *DefaultRequest) SetRequest(i interface{}) {
	dr.req = i
}
