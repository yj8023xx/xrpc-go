package main

import (
	"tinyrpc/server"
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
	s.AddService("HelloService", new(HelloService))
	s.Start()
}
