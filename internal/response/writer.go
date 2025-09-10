package response

import (
	"fmt"
	"io"

	"github.com/xixotron/httpfromtcp/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateWriteHeaders
	writerStateWriteBody
	writerStateDone
)

type Writer struct {
	state  writerState
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state:  writerStateStatusLine,
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.state)
	}
	defer func() { w.state = writerStateWriteHeaders }()

	_, err := fmt.Fprint(w.writer, getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateWriteHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.state)
	}
	defer func() { w.state = writerStateWriteBody }()

	for key, value := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateWriteBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	defer func() { w.state = writerStateDone }()

	return w.writer.Write(p)
}
