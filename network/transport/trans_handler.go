package transport

import (
	"xrpc/network/protocol"
)

type TransHandler interface {
	Handle(message *protocol.Message) (*protocol.Message, error)
}
