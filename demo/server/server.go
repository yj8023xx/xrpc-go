package main

import (
	"xrpc/server"
)

type Args struct {
	Name string
}

type HelloService struct{}

func (h *HelloService) Hello(args *Args, reply *string) error {
	*reply = "Hello " + args.Name
	return nil
}

func main() {
	s, _ := server.NewServer("8090", "zookeeper://127.0.0.1:2181")
	s.AddService("com.smallc.xrpc.api.hello.HelloService", new(HelloService))
	s.Start()
}
