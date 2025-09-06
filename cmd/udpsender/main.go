package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal("Couldn't resolve udp addres:", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Couldn't start udp connection:", err)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Ctrl+C to exit.\n", address)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %s\n", err)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatalf("Error writing to remote: %s\n", err)
		}
		fmt.Printf("Message sent: %s", message)
	}
}
