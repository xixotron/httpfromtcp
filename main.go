package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "./messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Couldn't open file %s: %s\n", inputFilePath, err)
	}
	defer f.Close()

	fmt.Printf("Reading file %s\n", inputFilePath)

	for {
		cont := make([]byte, 8)
		_, err := f.Read(cont)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		str := string(cont)
		fmt.Printf("read: %s\n", str)
	}
}
