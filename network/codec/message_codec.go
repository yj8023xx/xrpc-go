package codec

import (
	"bytes"
	"io"
	"xrpc/network/protocol"
)

type MessageCodec interface {
	ReadHeader(io.ReadWriteCloser, *protocol.Header) error
	Write(io.ReadWriteCloser, *protocol.Message) error
}

type messageCodec struct {
	headerCodec HeaderCodec
}

func (c *messageCodec) ReadHeader(conn io.ReadWriteCloser, header *protocol.Header) error {
	headerBytes := make([]byte, protocol.FixedHeaderLength)
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return err
	}
	return c.headerCodec.Decode(headerBytes, header)
}

func (c *messageCodec) Write(conn io.ReadWriteCloser, message *protocol.Message) error {
	headerBytes, err := c.headerCodec.Encode(message.Header)
	if err != nil {
		return err
	}
	byteArr := [][]byte{headerBytes, message.Payload}
	combineBytes := bytes.Join(byteArr, []byte(""))
	_, err = conn.Write(combineBytes)
	if err != nil {
		return err
	}
	return nil
}

var m MessageCodec = &messageCodec{
	headerCodec: NewHeaderCodec(),
}

func NewMessageCodec() MessageCodec {
	return m
}
