package server

import (
	"reflect"
	"strings"
	"xrpc/network/codec"
	"xrpc/network/protocol"
	"xrpc/network/transport"
)

type RpcRequestHandler struct {
	serviceMap map[string]*service
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (h *RpcRequestHandler) Handle(message *protocol.Message) (*protocol.Message, error) {
	header := message.Header
	name := codec.GetPayloadCodecName(int(header.SerializationId))
	payloadCodec := codec.GetPayloadCodec(name)
	rpcRequest := &protocol.RpcRequestPayload{}
	payloadCodec.Unmarshal(message.Payload, rpcRequest)
	service := h.serviceMap[rpcRequest.ServiceName]
	serviceMethod := service.methodMap[firstUpper(rpcRequest.MethodName)]
	argv := reflect.New(serviceMethod.argsType)
	kvPairs := rpcRequest.ArgMap
	for name, value := range kvPairs {
		argv.Elem().FieldByName(firstUpper(name)).Set(reflect.ValueOf(value))
	}
	replv := reflect.New(serviceMethod.replyType)
	service.call(serviceMethod, argv, replv)
	payload, err := payloadCodec.Marshal(replv.Interface())
	if err != nil {
		return nil, err
	}
	respHeader := &protocol.Header{
		MagicNumber:     protocol.MagicNumber,
		Version:         protocol.Version,
		HeaderLength:    uint16(protocol.FixedHeaderLength),
		TotalLength:     uint16(protocol.FixedHeaderLength + len(payload)),
		MessageTypeId:   protocol.RpcResponse,
		SerializationId: header.SerializationId,
		RequestId:       header.RequestId,
	}
	return &protocol.Message{
		Header:  respHeader,
		Payload: payload,
	}, nil
}

func NewRpcRequestHandler(serviceMap map[string]*service) transport.TransHandler {
	return &RpcRequestHandler{
		serviceMap: serviceMap,
	}
}
