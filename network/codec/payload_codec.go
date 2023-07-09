package codec

import (
	"encoding/json"
)

const (
	JsonName    string = "json"
	HessianName string = "hessian"
)

const (
	Json = iota
	Hessian
)

var payloadMap = make(map[string]PayloadCodec)
var nameMap = make(map[int]string)
var idMap = make(map[string]int)

func init() {
	payloadMap[JsonName] = NewJsonCodec()
	nameMap[Json] = JsonName
	idMap[JsonName] = Json
}

func GetPayloadCodecName(id int) string {
	return nameMap[id]
}

func GetPayloadCodecId(name string) int {
	return idMap[name]
}

// GetPayloadCodec gets desired payload codec from message.
func GetPayloadCodec(name string) PayloadCodec {
	pc := payloadMap[name]
	return pc
}

// PayloadCodec is used to marshal and unmarshal payload.
type PayloadCodec interface {
	Marshal(i interface{}) ([]byte, error)
	Unmarshal(data []byte, i interface{}) error
	Name() string
}

type JsonCodec struct{}

func (_ *JsonCodec) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (_ *JsonCodec) Unmarshal(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func (_ *JsonCodec) Name() string {
	return "json"
}

var jsonCodec PayloadCodec = &JsonCodec{}

func NewJsonCodec() PayloadCodec {
	return jsonCodec
}
