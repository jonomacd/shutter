package shutter

import (
	//"time"
	"fmt"

	message "github.com/jonomacd/shutter/proto"

	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	//"golang.org/x/net/context"
	"gopkg.in/project-iris/iris-go.v1"
)

var (
	ServiceRegistry Registry
	DefaultPort     int = 55555
)

type Registry interface {
	Name() string
	Register(string, HandlerFunc) error
	GetEndpoint(string) Endpoint
	iris.ServiceHandler
}

type Endpoint interface {
	Name() string
	Handler() HandlerFunc
}

// TODO Change return type to some kind of response object
type HandlerFunc func(req Request) ([]byte, error)

type DefaultEndpoint struct {
	name    string
	handler HandlerFunc
}

func (de *DefaultEndpoint) Name() string {
	return de.name
}

func (de *DefaultEndpoint) Handler() HandlerFunc {
	return de.handler
}

// Implements iris.ServiceHandler
type DefaultRegistry struct {
	service   *iris.Service
	name      string
	endpoints map[string]Endpoint
}

func InitializeService(name string) error {
	dr := &DefaultRegistry{
		name:      name,
		endpoints: make(map[string]Endpoint),
	}

	// TODO Specify port
	service, err := iris.Register(DefaultPort, name, dr, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	dr.service = service
	ServiceRegistry = dr

	return err

}

func (dr *DefaultRegistry) Name() string {
	return dr.name
}

func (dr *DefaultRegistry) Register(name string, handler HandlerFunc) error {
	dr.endpoints[name] = &DefaultEndpoint{
		name:    name,
		handler: handler,
	}
	return nil
}

func (dr *DefaultRegistry) GetEndpoint(name string) Endpoint {
	return dr.endpoints[name]
}

// Called once after the service is registered in the Iris network, but before
// and handlers are activated. Its goal is to initialize any internal state
// dependent on the connection.
func (dr *DefaultRegistry) Init(conn *iris.Connection) error {
	return nil
}

// Callback invoked whenever a broadcast message arrives designated to the
// cluster of which this particular service instance is part of.
func (dr *DefaultRegistry) HandleBroadcast(message []byte) {
	// Unmarshal the message and check the header
	// run the handler for that specific endpoint
	log.Info("broadcast")
}

// Callback invoked whenever a request designated to the service's cluster is
// load-balanced to this particular service instance.
//
// The method should service the request and return either a reply or the
// error encountered, which will be delivered to the request originator.
//
// Returning nil for both or none of the results will result in a panic. Also,
// since the requests cross language boundaries, only the error string gets
// delivered remotely (any associated type information is effectively lost).
func (dr *DefaultRegistry) HandleRequest(request []byte) ([]byte, error) {

	in := &message.Request{}
	err := proto.Unmarshal(request, in)
	if err != nil {
		// Figure out how to marshal the error response
		return nil, err
	}

	handlerRequest := NewRequest(in.Originator, in.Endpoint, in.Body, in.Headers)

	// Set the clientside timeout so that handlers can give up in situations they are running long
	/*to, _ := time.ParseDuration(in.ClientTimeout)
	if to != 0 {
		handlerRequest.Context, cancel := context.WithTimeout(context.Background(), to)
		defer cancel()
	}*/

	ep := ServiceRegistry.GetEndpoint(in.Endpoint)
	if ep == nil {
		// construct endpoint not found response
		return nil, fmt.Errorf("No Ep!")
	}
	rsp, err := ep.Handler()(handlerRequest)
	// marshal the response and return the reply or the error

	return rsp, err
}

// Callback invoked whenever a tunnel designated to the service's cluster is
// constructed from a remote node to this particular instance.
func (dr *DefaultRegistry) HandleTunnel(tunnel *iris.Tunnel) {
	// Todo decide if I want to handle this
	log.Info("tunnel")
}

// Callback notifying the service that the local relay dropped its connection.
func (dr *DefaultRegistry) HandleDrop(reason error) {
	// Need to do error handling here
	log.Fatalf("[SHUTTER] Disconnected from relay: %v", reason)

}
