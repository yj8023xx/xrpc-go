package net

import (
	"io"
	"net"
	"sync"
	"time"
	"tinyrpc/network/codec"
	"tinyrpc/network/protocol"
	"tinyrpc/network/transport"
)

type netClient struct {
	conns []net.Conn
}

func (c *netClient) CreateTransport(address string, timeout time.Duration) (transport.Transport, error) {
	conn, _ := net.DialTimeout("tcp", address, timeout)
	c.conns = append(c.conns, conn)
	return newNetTransport(conn), nil
}

func (c *netClient) Close() error {
	for _, conn := range c.conns {
		conn.Close()
	}
	return nil
}

type netTransport struct {
	mu            sync.Mutex
	conn          io.ReadWriteCloser
	messageCodec  codec.MessageCodec
	inFlightCalls transport.InFlightCalls
}

func (t *netTransport) Send(message *protocol.Message, reply interface{}) (*transport.Call, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.messageCodec.Write(t.conn, message)
	call := &transport.Call{
		RequestId: int(message.Header.RequestId),
		Reply:     reply,
		Done:      make(chan *transport.Call, 1),
	}
	t.inFlightCalls.AddCall(call)
	return call, nil
}

func (t *netTransport) receive() error {
	var header protocol.Header
	t.messageCodec.ReadHeader(t.conn, &header)
	payload := make([]byte, header.TotalLength-header.HeaderLength)
	io.ReadFull(t.conn, payload)
	call, _ := t.inFlightCalls.RemoveCall(int(header.RequestId))
	// 根据序列化id反序列化
	name := codec.GetPayloadCodecName(int(header.SerializationId))
	payloadCodec := codec.GetPayloadCodec(name)
	payloadCodec.Unmarshal(payload, call.Reply)
	call.Complete()
	return nil
}

func NewNetClient() transport.TransClient {
	return &netClient{
		conns: []net.Conn{},
	}
}

func newNetTransport(conn net.Conn) transport.Transport {
	transport := &netTransport{
		mu:            sync.Mutex{},
		conn:          conn,
		messageCodec:  codec.NewMessageCodec(),
		inFlightCalls: transport.InFlightCalls{},
	}
	go transport.receive()
	return transport
}
