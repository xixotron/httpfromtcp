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

const httpVersion = "HTTP/1.1"

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("trying to write status line in incorrect state %d", w.state)
	}

	_, err := fmt.Fprintf(w.writer, "%s %d %s\r\n", httpVersion, statusCode, statusCodeText(statusCode))
	if err != nil {
		return err
	}

	w.state = writerStateWriteHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateWriteHeaders {
		return fmt.Errorf("trying to write headers in incorrect state %d", w.state)
	}

	for key, value := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.writer, "\r\n")
	if err != nil {
		return err
	}

	w.state = writerStateWriteBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateWriteBody {
		return 0, fmt.Errorf("trying to write body in incorrect state %d", w.state)
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}

	w.state = writerStateDone
	return n, nil
}
