package response

import (
	"fmt"
	"io"

	"github.com/xixotron/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

const httpVersion = "HTTP/1.1"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := fmt.Fprintf(w, "%s %d %s\r\n", httpVersion, statusCode, statusCodeText(statusCode))
	return err
}

func statusCodeText(statusCode StatusCode) (statusText string) {
	switch statusCode {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "\r\n")
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", 10))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}
