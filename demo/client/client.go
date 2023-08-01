package main

import (
	"fmt"
	"xrpc/client"
	"xrpc/network/codec"
	"xrpc/registry"
)

type Args struct {
	Name string
}

func main() {
	r, _ := registry.NewRegistry("zookeeper://127.0.0.1:2181")
	xClient := client.NewXClient("com.smallc.xrpc.api.hello.HelloService", r, client.RoundRobin, codec.Json)
	args := &Args{Name: "World"}
	var reply string
	err := xClient.Call("Hello", args, &reply)
	if err != nil {

	}
	fmt.Printf("reply: %s\n", reply)
}
