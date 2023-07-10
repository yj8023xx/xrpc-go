package protocol

const (
	FixedHeaderLength int = 17
)

const (
	RpcRequest = iota
	RpcResponse
)

const (
	MagicNumber = 0x8023
	Version     = 1
)

// Header 12 bytes
type Header struct {
	MagicNumber     uint16
	Version         uint8
	HeaderLength    uint16
	TotalLength     uint16
	MessageTypeId   uint8
	SerializationId uint8
	RequestId       uint64
}

type RpcRequestPayload struct {
	ServiceName string
	MethodName  string
	Args        interface{}
}

type Message struct {
	*Header
	Payload []byte
}
