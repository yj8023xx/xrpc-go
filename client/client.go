package client

import (
	"sync"
	"time"
	"tinyrpc/network/transport"
	_net "tinyrpc/network/transport/net"
)

type Client interface {
	CreateTransport(uri string) (transport.Transport, error)
	Close() error
}

type client struct {
	mu           sync.Mutex
	transportMap map[string]transport.Transport
	transClient  transport.TransClient
}

func (c *client) Close() error {
	return c.transClient.Close()
}

func newClient() Client {
	return &client{
		mu:           sync.Mutex{},
		transportMap: make(map[string]transport.Transport),
		transClient:  _net.NewNetClient(),
	}
}

func (c *client) CreateTransport(uri string) (transport.Transport, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.transportMap[uri] == nil {
		var err error
		c.transportMap[uri], err = c.transClient.CreateTransport(uri, time.Second)
		if err != nil {
			return nil, err
		}
	}
	return c.transportMap[uri], nil
}

var c = newClient()
