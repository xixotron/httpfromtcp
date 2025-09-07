package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}

	str := strings.TrimSpace(string(data[:idx]))
	key, value, err := headerFromString(str)
	if err != nil {
		return 0, false, err
	}
	h.Set(key, value)

	return idx + len(crlf), false, nil
}

func headerFromString(str string) (key string, value string, err error) {
	parts := strings.SplitN(str, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header line %s\n", str)
	}

	key = parts[0]
	if key != strings.TrimRight(key, " ") {
		return "", "", fmt.Errorf("invalid header name %s\n", key)
	}

	value = strings.TrimSpace(parts[1])
	return key, value, nil
}

func (h Headers) Set(key string, value string) {
	h[key] = value
}
