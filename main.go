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
	lineChan := make(chan string)
	go func() {
		defer f.Close()
		defer close(lineChan)

		line := ""

		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			if err != nil {
				if line != "" {
					lineChan <- line
				}
				if errors.Is(err, io.EOF) {
					break
				}

				log.Printf("Error reading file: %v\n", err)
				break
			}

			parts := strings.Split(string(buff[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				lineChan <- line + parts[i]
				line = ""
			}
			line += parts[len(parts)-1]

		}
	}()

	return lineChan
}
