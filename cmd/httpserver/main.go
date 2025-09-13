package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/xixotron/httpfromtcp/internal/headers"
	"github.com/xixotron/httpfromtcp/internal/request"
	"github.com/xixotron/httpfromtcp/internal/response"
	"github.com/xixotron/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerFunc(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		req.RequestLine.RequestTarget = strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
		handleHTTPProxy(w, req, "https://httpbin.org")
	} else if req.RequestLine.RequestTarget == "/video" {
		handleVideo(w, req)
	} else {
		handler200(w, req)
	}
}

const template = `<html>
  <head>
    <title>{{title}}</title>
  </head>
  <body>
    <h1>{{heading}}</h1>
    <p>{{paragraph}}</p>
  </body>
</html>
`

func replaceTemplate(template, title, heading, paragraph string) string {
	replacer := strings.NewReplacer(
		"{{title}}", title,
		"{{heading}}", heading,
		"{{paragraph}}", paragraph,
	)
	return replacer.Replace(template)

}

func handler400(w *response.Writer, _ *request.Request) {
	resp := replaceTemplate(template,
		"400 Bad Request",
		"Bad Request",
		"Your request honestly kinda sucked.",
	)
	headers := response.GetDefaultHeaders(len(resp))
	headers.Override("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(resp))
}

func handler500(w *response.Writer, _ *request.Request) {
	resp := replaceTemplate(template,
		"500 Internal Server Error",
		"Internal Server Error",
		"Okay, you know what? This one is on me.",
	)
	headers := response.GetDefaultHeaders(len(resp))
	headers.Override("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(resp))
}

func handler200(w *response.Writer, _ *request.Request) {
	resp := replaceTemplate(template,
		"200 OK",
		"Success!",
		"Your request was an absolute banger.",
	)
	headers := response.GetDefaultHeaders(len(resp))
	headers.Override("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(resp))
}

func handleHTTPProxy(w *response.Writer, req *request.Request, target string) {
	url, err := url.JoinPath(target, req.RequestLine.RequestTarget)
	if err != nil {
		handler400(w, req)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error redirecting request: %v", err)
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
	h.Override("Content-Type", resp.Header.Get("Content-Type"))
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(h)

	const chunkSize = 1024
	buff := make([]byte, chunkSize)
	hash := sha256.New()
	contentLength := 0
	for {
		n, err := resp.Body.Read(buff)
		log.Printf("Reading %v bytes from target", n)
		if n > 0 {
			contentLength += n
			_, err := w.WriteChunkedBody(buff[:n])
			hash.Write(buff[:n])
			if err != nil {
				log.Printf("error writing chunked body: %v", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading response body: %v", err)
			break
		}
	}

	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash.Sum(nil)))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", contentLength))
	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Printf("error writing trailers: %v", err)
	}
}

func handleVideo(w *response.Writer, req *request.Request) {
	file, err := os.Open("./assets/vim.mp4")
	if err != nil {
		handler500(w, req)
		log.Printf("error reading video file: %v", err)
		return
	}

	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
	h.Override("Content-Type", "video/mp4")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(h)

	const chunkSize = 1024
	buff := make([]byte, chunkSize)
	for {
		n, err := file.Read(buff)
		log.Printf("Reading %v video bytes", n)
		if n > 0 {
			_, err := w.WriteChunkedBody(buff[:n])
			if err != nil {
				log.Printf("error writing chunked body: %v", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading response body: %v", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("error writing trailers: %v", err)
	}
}
