package client

import (
	"github.com/bytedance/gopkg/lang/fastrand"
	"xrpc/network/transport"
)

const (
	Random = iota
	RoundRobin
)

type Selector interface {
	Select() transport.Transport
	Update([]transport.Transport)
}

type RandomSelector struct {
	transports []transport.Transport
}

func (s *RandomSelector) Select() transport.Transport {
	if len(s.transports) == 0 {
		return nil
	}
	i := fastrand.Intn(len(s.transports))
	return s.transports[i]
}

func (s *RandomSelector) Update(transports []transport.Transport) {
	s.transports = transports
}

func newRandomSelector(transports []transport.Transport) Selector {
	return &RandomSelector{
		transports: transports,
	}
}

type RoundRobinSelector struct {
	transports []transport.Transport
	i          int
}

func (s *RoundRobinSelector) Select() transport.Transport {
	if len(s.transports) == 0 {
		return nil
	}
	i := s.i
	i = i % len(s.transports)
	s.i = i + 1
	return s.transports[i]
}

func (s *RoundRobinSelector) Update(transports []transport.Transport) {
	s.transports = transports
}

func newRoundRobinSelector(transports []transport.Transport) Selector {
	return &RoundRobinSelector{
		transports: transports,
		i:          0,
	}
}

type ConsistentHashSelector struct {
}

func newSelector(selectMode int, transports []transport.Transport) Selector {
	switch selectMode {
	case Random:
		return newRandomSelector(transports)
	case RoundRobin:
		return newRoundRobinSelector(transports)
	default:
		return newRandomSelector(transports)
	}
}
