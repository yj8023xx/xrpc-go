package main

import (
	"fmt"
	"tinyrpc/client"
	"tinyrpc/network/codec"
	"tinyrpc/registry"
)

type Args struct {
	Name string
}

func main() {
	r, _ := registry.GetRegistry("zookeeper://127.0.0.1:2181")
	xClient := client.NewXClient("HelloService", r, client.RoundRobin, codec.Json)
	args := &Args{Name: "World"}
	var reply string
	err := xClient.Call("Hello", args, &reply)
	if err != nil {

	}
	fmt.Printf("reply: %s\n", reply)
}
