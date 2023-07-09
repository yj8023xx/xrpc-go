package transport

import (
	"sync"
	"time"
	"tinyrpc/network/protocol"
)

type TransClient interface {
	CreateTransport(address string, timeout time.Duration) (Transport, error)
	Close() error
}

type Transport interface {
	Send(*protocol.Message, interface{}) (*Call, error)
}

type Call struct {
	RequestId int
	Reply     interface{} // The reply from the function (*struct).
	Done      chan *Call
}

func (call *Call) Complete() {
	call.Done <- call
}

type InFlightCalls struct {
	callMap sync.Map
}

func (i *InFlightCalls) AddCall(call *Call) error {
	i.callMap.Store(call.RequestId, call)
	return nil
}

func (i *InFlightCalls) RemoveCall(requestId int) (*Call, error) {
	call, _ := i.callMap.Load(requestId)
	i.callMap.Delete(requestId)
	return call.(*Call), nil
}
