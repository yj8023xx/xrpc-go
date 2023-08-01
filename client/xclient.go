package client

import (
	"github.com/google/uuid"
	"net/url"
	"reflect"
	"strings"
	"xrpc/network/codec"
	"xrpc/network/protocol"
	"xrpc/network/transport"
	"xrpc/registry"
)

type XClient interface {
	Go(methodName string, args, reply interface{}) *transport.Call
	Call(methodName string, args, reply interface{}) error
}

type xClient struct {
	serviceName  string
	registry     registry.Registry
	selector     Selector
	payloadCodec codec.PayloadCodec
}

func firstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func convertToMap(s interface{}) map[string]interface{} {
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	m := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		m[firstLower(field.Name)] = value
	}
	return m
}

// Go invokes the function asynchronously. It returns the Call structure representing the invocation.
// The done channel will signal when the call is complete by returning the same Call object.
func (c *xClient) Go(methodName string, args, reply interface{}) *transport.Call {
	rpcRequest := protocol.RpcRequestPayload{
		ServiceName: c.serviceName,
		MethodName:  firstLower(methodName),
		ArgMap:      convertToMap(args),
	}
	payload, _ := c.payloadCodec.Marshal(rpcRequest)
	uuid, _ := uuid.NewUUID()
	header := &protocol.Header{
		MagicNumber:     protocol.MagicNumber,
		Version:         protocol.Version,
		HeaderLength:    uint16(protocol.FixedHeaderLength),
		TotalLength:     uint16(protocol.FixedHeaderLength + len(payload)),
		MessageTypeId:   protocol.RpcRequest,
		SerializationId: uint8(codec.GetPayloadCodecId(c.payloadCodec.Name())),
		RequestId:       uint64(uuid.ID()),
	}
	message := &protocol.Message{
		Header:  header,
		Payload: payload,
	}
	call, _ := c.selector.Select().Send(message, reply)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *xClient) Call(methodName string, args, reply interface{}) error {
	<-c.Go(methodName, args, reply).Done
	return nil
}

// watch changes of service and update transports.
func (c *xClient) watch() {
	var transports []transport.Transport
	c.selector.Update(transports)
}

func NewXClient(serviceName string, registry registry.Registry, selectMode int, serializationMethod int) XClient {
	uris, err := registry.GetServiceAddress(serviceName)
	if err != nil {

	}
	var transports []transport.Transport
	for _, uri := range uris {
		u, _ := url.Parse(uri)
		transport, _ := c.CreateTransport(u.Host)
		transports = append(transports, transport)
	}
	return &xClient{
		serviceName:  serviceName,
		registry:     registry,
		selector:     newSelector(selectMode, transports),
		payloadCodec: codec.GetPayloadCodec(codec.GetPayloadCodecName(serializationMethod)),
	}
}
