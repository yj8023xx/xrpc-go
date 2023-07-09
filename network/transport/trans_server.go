package transport

type TransServerFactory interface {
	NewTransServer() TransServer
}

type TransServer interface {
	Start(address string, handlerMap map[int]TransHandler) error
}
