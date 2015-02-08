package server

type Endpoint interface {
	Name() string
	Handler() HandlerFunc
	Request() interface{}
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
	return de.handler
}

func (de *DefaultEndpoint) Request() interface{} {
	return de.req
}
