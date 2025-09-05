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
	defer f.Close()

	fmt.Printf("Reading file %s\n", inputFilePath)

	line := ""
	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			break
		}

		parts := strings.Split(string(b[:n]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", line+parts[i])
			line = ""
		}
		line += parts[len(parts)-1]
	}

	if line != "" {
		fmt.Printf("read: %q\n", line)
		line = ""
	}
}
