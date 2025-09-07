package main

import (
	"fmt"
	"log"
	"net"

	"github.com/xixotron/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Couldn't start a socket on port %s/tcp: %s\n", port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s/tcp\n", port)
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s\n", err)
		}
		fmt.Printf("Accepted connection from %s\n", con.RemoteAddr())

		req, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatalf("Error parsing request: %s\n", err)
		}
		printRequest(req)

		fmt.Printf("connection from %s has been closed\n", con.RemoteAddr())
	}
}

func printRequest(req *request.Request) {
	fmt.Println("Request line:")
	fmt.Println("- Method:", req.RequestLine.Method)
	fmt.Println("- Target:", req.RequestLine.RequestTarget)
	fmt.Println("- Version:", req.RequestLine.HttpVersion)
}
