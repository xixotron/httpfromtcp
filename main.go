package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		for line := range getLinesChannel(con) {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Printf("connection from %s has been closed\n", con.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)

		line := ""
		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Printf("Error reading file: %v\n", err)
				}
				break
			}

			str := string(buff[:n])
			parts := strings.Split(str, "\n")
			for i, part := range parts {
				if i > 0 {
					lines <- line
					line = ""
				}
				line += part
			}
		}
		if line != "" {
			lines <- line
		}
	}()

	return lines
}
