package codec

import (
	"xrpc/network/protocol"
)

type HeaderCodec interface {
	Encode(header *protocol.Header) ([]byte, error)
	Decode([]byte, *protocol.Header) error
}

type headerCodec struct {
}

func (c *headerCodec) Encode(header *protocol.Header) ([]byte, error) {
	byteBuffer := NewWriterBuffer(protocol.FixedHeaderLength)
	WriteUint16(header.MagicNumber, byteBuffer)
	WriteByte(header.Version, byteBuffer)
	WriteUint16(header.HeaderLength, byteBuffer)
	WriteUint16(header.TotalLength, byteBuffer)
	WriteByte(header.MessageTypeId, byteBuffer)
	WriteByte(header.SerializationId, byteBuffer)
	WriteUint64(header.RequestId, byteBuffer)
	WriteByte(header.Status, byteBuffer)
	return byteBuffer.Bytes()
}

func (c *headerCodec) Decode(in []byte, header *protocol.Header) error {
	byteBuffer := NewReaderBuffer(in)
	header.MagicNumber, _ = ReadUint16(byteBuffer)
	header.Version, _ = ReadByte(byteBuffer)
	header.HeaderLength, _ = ReadUint16(byteBuffer)
	header.TotalLength, _ = ReadUint16(byteBuffer)
	header.MessageTypeId, _ = ReadByte(byteBuffer)
	header.SerializationId, _ = ReadByte(byteBuffer)
	header.RequestId, _ = ReadUint64(byteBuffer)
	header.Status, _ = ReadByte(byteBuffer)
	return nil
}

var h HeaderCodec = &headerCodec{}

func NewHeaderCodec() HeaderCodec {
	return h
}
