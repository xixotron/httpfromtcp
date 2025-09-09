package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xixotron/httpfromtcp/internal/request"
	"github.com/xixotron/httpfromtcp/internal/response"
	"github.com/xixotron/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerFunc(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your problem is not my problem\n",
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}

	_, err := w.Write([]byte("All good, frfr\n"))
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
	return nil
}
