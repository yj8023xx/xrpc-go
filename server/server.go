package server

import (
	"net"
	"tinyrpc/network/protocol"
	"tinyrpc/network/transport"
	_net "tinyrpc/network/transport/net"
	"tinyrpc/registry"
)

type Server interface {
	Start()
	AddService(string, interface{})
}

type server struct {
	host        string
	port        string
	uri         string
	serviceMap  map[string]*service
	handlerMap  map[int]transport.TransHandler
	registry    registry.Registry
	transServer transport.TransServer
}

func (s *server) Start() {
	if s.transServer == nil {
		s.transServer = _net.NewNetServer()
		s.handlerMap[protocol.RpcRequest] = NewRpcRequestHandler(s.serviceMap)
		s.transServer.Start(s.host+":"+s.port, s.handlerMap)
	}
}

func (s *server) AddService(serviceName string, rcvr interface{}) {
	service := NewService(rcvr)
	s.serviceMap[serviceName] = service
	s.registry.RegisterService(serviceName, s.uri)
}

func NewServer(port string, nameServiceUri string) (Server, error) {
	registry, err := registry.GetRegistry(nameServiceUri)
	if err != nil {
		return nil, err
	}

	// get local ip address
	interfaces, err := net.Interfaces()
	var host string
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue
		}
		if (i.Flags & net.FlagLoopback) != 0 {
			continue
		}
		addresses, _ := i.Addrs()
		for _, addr := range addresses {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					host = ipNet.IP.String()
					break
				}
			}
		}
	}
	return &server{
		host:       host,
		port:       port,
		uri:        "rpc://" + host + ":" + port,
		serviceMap: make(map[string]*service),
		handlerMap: make(map[int]transport.TransHandler),
		registry:   registry,
	}, nil
}
