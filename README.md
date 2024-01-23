# xRPC-Go

[xRPC](https://github.com/yj8023xx/xrpc) has been implemented in Go, providing a lightweight, high-throughput, low-latency RPC framework.



## Example

**Client**

```go
type Args struct {
	Name string
}

func main() {
	r, _ := registry.NewRegistry("zookeeper://127.0.0.1:2181")
	xClient := client.NewXClient("HelloService", r, client.RoundRobin, codec.Json)
	args := &Args{Name: "World"}
	var reply string
	err := xClient.Call("Hello", args, &reply)
	if err != nil {

	}
	fmt.Printf("reply: %s\n", reply)
}
```

**Server**

```go
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
```
