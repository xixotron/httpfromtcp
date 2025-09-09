package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xixotron/httpfromtcp/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	state       requestState
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

const bufferSize = 8

func RequestFromReader(r io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
	for request.state != requestStateDone {
		if readToIndex >= len(buff) {
			tmp := make([]byte, len(buff)*2)
			copy(tmp, buff)
			buff = tmp
		}

		bytesRead, err := r.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.state, bytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += bytesRead

		bytesParsed, err := request.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[bytesParsed:])
		readToIndex -= bytesParsed
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requetLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requetLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		bodyLenStr := r.Headers.Get("Content-Length")
		if bodyLenStr == "" {
			r.state = requestStateDone
			return 0, nil
		}
		bodyLength, err := strconv.Atoi(bodyLenStr)
		if err != nil {
			return 0, fmt.Errorf("error: parsing Content-Length: %v", err)
		}
		r.Body = append(r.Body, data...)
		if len(r.Body) > bodyLength {
			return 0, fmt.Errorf("error: body longer that Content-Length")
		}
		if len(r.Body) == bodyLength {
			r.state = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + len(crlf), nil
}

func requestLineFromString(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Malformed request line")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("Invalid Method")
		}
	}

	requestTarget := parts[1]

	httpVersionParts := strings.Split(parts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, fmt.Errorf("Malformed HTTP-Version")
	}

	httpPart := httpVersionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("Unknown HTTP-Version")
	}
	version := httpVersionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("Unknown / unsupported HTTP-version")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
