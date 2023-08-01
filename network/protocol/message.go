package protocol

const (
	RpcRequest = iota
	RpcResponse
)

const (
	MagicNumber       = 0x8023
	Version           = 1
	FixedHeaderLength = 18
)

// Header 18 bytes
type Header struct {
	MagicNumber     uint16
	Version         uint8
	HeaderLength    uint16
	TotalLength     uint16
	MessageTypeId   uint8
	SerializationId uint8
	RequestId       uint64
	Status          uint8
}

type RpcRequestPayload struct {
	ServiceName string                 `json:"serviceName"`
	MethodName  string                 `json:"methodName"`
	ArgMap      map[string]interface{} `json:"argMap"`
}

type Message struct {
	*Header
	Payload []byte
}
