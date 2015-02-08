package client

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	message "github.com/jonomacd/shutter/proto"
	"gopkg.in/project-iris/iris-go.v1"
)

var (
	Port         int = 55555
	GlobalClient Client
)

type Client interface {
	// Allows for setting functions to call prior to the request.
	// The idea here is that we can hook in rate limiting, circuit breaking, trace, etc
	// If an error is returned then the request is canceled before it even starts
	SetPreRequest(func(service, endpoint string) error)

	// Make a request with default options
	Request(service, endpoint string, req, rsp proto.Message, opts ...*Options) error

	// Set default options for this client
	SetOptions(options *Options)

	// Allows for setting functions to call after the request. Required for circuit breaking,
	// trace, etc
	SetPostRequest(func(service, endpoint string, dur time.Duration, err error))
}

// TODO: This will include things like timeout, retry policy, etc
type Options struct {
	Timeout time.Duration
}

type DefaultClient struct {
	conn             *iris.Connection
	defaultOptions   *Options
	preRequestHooks  []func(service, endpoint string) error
	postRequestHooks []func(service, endpoint string, dur time.Duration, err error)
}

// Builds a global client
func InitializeClient() error {
	var err error
	GlobalClient, err = NewClient(Port, nil)
	return err
}

// Request sends a request via the GlobalClient
func Request(service, endpoint string, req, rsp proto.Message, options ...*Options) error {
	if GlobalClient == nil {
		// Alternatively we could just do the init for you in this case
		// and not return an error.  Might make more sense though it is
		// nice to make people think about using a bespoke client.
		// InitializeClient()
		return fmt.Errorf("Global Client Not Initialized")
	}
	return GlobalClient.Request(service, endpoint, req, rsp, options...)
}

// NewClient builds a new client and with it a new connection to iris.
func NewClient(port int, options *Options) (Client, error) {
	conn, err := iris.Connect(port)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &Options{
			Timeout: time.Second * 5,
		}
	}

	return &DefaultClient{
		conn:           conn,
		defaultOptions: options,
	}, nil
}

// Request sends a request over iris.  It will send the req and fill in the rsp object
func (dc *DefaultClient) Request(service, endpoint string, req, rsp proto.Message, opts ...*Options) error {
	if dc == nil {
		return fmt.Errorf("Client Not Initialized")
	}

	for _, hook := range dc.preRequestHooks {
		err := hook(service, endpoint)
		if err != nil {
			return err
		}
	}

	// Set the options for the request
	options := dc.defaultOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	wireReq := &message.Request{
		Endpoint:      endpoint,
		Originator:    "",
		ClientTimeout: options.Timeout.String(),
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	wireReq.Body = b

	// Todo error checking
	wireb, _ := proto.Marshal(wireReq)
	startTime := time.Now()
	rspb, err := dc.conn.Request(service, wireb, options.Timeout)
	dur := time.Since(startTime)
	for _, hook := range dc.postRequestHooks {
		defer hook(service, endpoint, dur, err)
	}
	if err != nil {
		return err
	}

	wireRsp := &message.Response{}

	proto.Unmarshal(rspb, wireRsp)

	if wireRsp.Type == "error" {
		// Todo, deal with errors better.  Build specific types based on the code
		err = fmt.Errorf("%s :: %s", wireRsp.Err.ErrorText, wireRsp.Err.Code)
		return err
	}

	err = proto.Unmarshal(wireRsp.Body, rsp)
	return err
}

func (dc *DefaultClient) SetOptions(options *Options) {}

func (dc *DefaultClient) SetPreRequest(f func(service, endpoint string) error) {
	dc.preRequestHooks = append(dc.preRequestHooks, f)
}
func (dc *DefaultClient) SetPostRequest(f func(service, endpoing string, dur time.Duration, err error)) {
	dc.postRequestHooks = append(dc.postRequestHooks, f)
}
