package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	state       requestState
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
		state: requestStateInitialized,
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
				request.state = requestStateDone
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
	switch r.state {
	case requestStateInitialized:
		requetLine, c, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if c == 0 {
			return 0, nil
		}
		r.RequestLine = *requetLine
		r.state = requestStateDone
		return c, nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
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
		return nil, errors.New("Malformed request line")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, errors.New("Invalid Method")
		}
	}

	requestTarget := parts[1]

	httpVersionParts := strings.Split(parts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, errors.New("Malformed HTTP-Version")
	}

	httpPart := httpVersionParts[0]
	if httpPart != "HTTP" {
		return nil, errors.New("Unknown HTTP-Version")
	}
	version := httpVersionParts[1]
	if version != "1.1" {
		return nil, errors.New("Unknown / unsupported HTTP-version")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
