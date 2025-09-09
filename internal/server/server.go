package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/xixotron/httpfromtcp/internal/request"
	"github.com/xixotron/httpfromtcp/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler:  handler,
		listener: listener,
	}
	go s.accept()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error Accepting connection : %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		herr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		herr.Write(conn)
		return
	}

	var buff bytes.Buffer
	herr := s.handler(&buff, req)
	if herr != nil {
		herr.Write(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("Error writing response Status-Line: %v", err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(buff.Len()))
	if err != nil {
		log.Printf("Error writing response Headers: %v", err)
	}

	_, err = buff.WriteTo(conn)
	if err != nil {
		log.Printf("Error writing response Body: %v", err)
	}
}
