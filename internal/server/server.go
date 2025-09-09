package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/xixotron/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, err
	}
	s := &Server{
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

const responseText = "HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n" +
	"Content-Length: 13\r\n" +
	"\r\n" +
	"Hello World!\n"

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("Error writing response Status-Line: %v\n", err)
		return
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Printf("Error writing response Headers: %v\n", err)
		return
	}
}
