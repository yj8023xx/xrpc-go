package net

import (
	"io"
	"net"
	"xrpc/network/codec"
	"xrpc/network/protocol"
	"xrpc/network/transport"
)

type netServer struct {
	messageCodec codec.MessageCodec
	handlerMap   map[int]transport.TransHandler
}

func (s *netServer) Start(address string, handlerMap map[int]transport.TransHandler) error {
	s.handlerMap = handlerMap
	listen, _ := net.Listen("tcp", address)
	for {
		conn, _ := listen.Accept()
		go s.handleRequest(conn)
	}
}

func (s *netServer) handleRequest(conn net.Conn) error {
	header := &protocol.Header{}
	s.messageCodec.ReadHeader(conn, header)
	payload := make([]byte, header.TotalLength-header.HeaderLength)
	_, err := io.ReadFull(conn, payload)
	if err != nil {
		return err
	}
	message := protocol.Message{
		Header:  header,
		Payload: payload,
	}
	// 获取对应的请求处理器
	handler := s.handlerMap[int(header.MessageTypeId)]
	// 提交给对应的请求处理器进行处理
	respMessage, _ := handler.Handle(&message)
	err = s.messageCodec.Write(conn, respMessage)
	if err != nil {
		return err
	}
	return nil
}

func NewNetServer() transport.TransServer {
	return &netServer{
		messageCodec: codec.NewMessageCodec(),
	}
}
