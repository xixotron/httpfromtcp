package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "./messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Couldn't open file %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading file %s\n", inputFilePath)

	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", line)
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
