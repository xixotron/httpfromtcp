package response

import (
	"fmt"
	"strings"

	"github.com/xixotron/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

const httpVersion = "HTTP/1.1"

func getStatusLine(statusCode StatusCode) (statusLine string) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %d ", httpVersion, statusCode))

	switch statusCode {
	case StatusOK:
		sb.WriteString("OK")
	case StatusBadRequest:
		sb.WriteString("Bad Request")
	case StatusInternalServerError:
		sb.WriteString("Internal Server Error")
	default:
		sb.WriteString("Unknown Status")
	}
	sb.WriteString("\r\n")
	return sb.String()
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}
