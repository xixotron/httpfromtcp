package server

import (
	"fmt"
	"io"
	"log"

	"github.com/xixotron/httpfromtcp/internal/request"
	"github.com/xixotron/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (herr *HandlerError) Write(w io.Writer) {
	if herr == nil {
		log.Println("call to handlerErrorReport without handlerError")
		return
	}

	errorMessage := fmt.Sprintf("Error: %s", herr.Message)

	err := response.WriteStatusLine(w, herr.StatusCode)
	if err != nil {
		log.Printf("Error writing handlerError Status-Line: %v", err)
	}

	err = response.WriteHeaders(w,
		response.GetDefaultHeaders(len(errorMessage)))
	if err != nil {
		log.Printf("Error writing handlerError Headers: %v", err)
	}

	_, err = w.Write([]byte(errorMessage))
	if err != nil {
		log.Printf("Error writing handlerError message: %v", err)
	}

	log.Printf("HandlerError: %d, %s",
		herr.StatusCode,
		herr.Message,
	)
}
