package server

import (
	"bufio"
	"net"
	"strconv"

	"github.com/Malayt04/redis-in-go/internal/logger"
	"github.com/Malayt04/redis-in-go/internal/resp"
	"github.com/Malayt04/redis-in-go/internal/handler"
	"github.com/Malayt04/redis-in-go/internal/store"
)

type Server struct {
	address string
	h       *handler.Handler
}

func New(address string) *Server {
	s := store.New()
	h := handler.New(s)
	return &Server{
		address: address,
		h:       h,
	}
}

func (s *Server) Start() {
	logger.Info("Starting Redis server on %s", s.address)
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}

	logger.Info("Server ready to accept connections")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Warn("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	logger.Info("Client connected: %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		msg, err := resp.Decode(reader)
		if err != nil {
			if err.Error() == "EOF" {
				logger.Debug("Client %s disconnected", conn.RemoteAddr())
			} else {
				logger.Error("Failed to decode message from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		logger.Debug("Received command: %v", msg)

		response := s.h.Handle(msg)

		if response.Type == handler.ResponseError {
			logger.Warn("Command error: %s", response.Value)
		}

		encoded := encodeResponse(response)

		_, err = conn.Write([]byte(encoded))
		if err != nil {
			logger.Error("Failed to write response to %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
}

func encodeResponse(r handler.Response) string {
	switch r.Type {
	case handler.ResponseSimpleString:
		return resp.EncodeString(r.Value)
	case handler.ResponseError:
		return resp.EncodeError(r.Value)
	case handler.ResponseBulkString:
		return resp.EncodeBulkString(r.Value)
	case handler.ResponseInteger:
		n, _ := strconv.Atoi(r.Value)
		return resp.EncodeInteger(n)
	case handler.ResponseNull:
		return resp.EncodeNull()
	default:
		return resp.EncodeBulkString(r.Value)
	}
}