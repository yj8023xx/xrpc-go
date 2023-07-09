package transport

import (
	"tinyrpc/network/protocol"
)

type TransHandler interface {
	Handle(message *protocol.Message) (*protocol.Message, error)
}
