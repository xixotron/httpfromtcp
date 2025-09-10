package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	var resp string
	var status response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		resp = replaceTemplate(template,
			"400 Bad Request",
			"Bad Request",
			"Your request honestly kinda sucked.",
		)
	case "/myproblem":
		resp = replaceTemplate(template,
			"500 Internal Server Error",
			"Internal Server Error",
			"Okay, you know what? This one is on me.",
		)
		status = response.StatusInternalServerError
	default:
		resp = replaceTemplate(template,
			"200 OK",
			"Success!",
			"Your request was an absolute banger.",
		)
		status = response.StatusOK
	}

	err := w.WriteStatusLine(status)
	if err != nil {
		log.Print(err)
	}

	headers := response.GetDefaultHeaders(len(resp))
	headers.Override("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		log.Print(err)
	}

	_, err = w.WriteBody([]byte(resp))
	if err != nil {
		log.Print(err)
	}
}

func replaceTemplate(template, title, heading, paragraph string) string {
	replacer := strings.NewReplacer(
		"{{title}}", title,
		"{{heading}}", heading,
		"{{paragraph}}", paragraph,
	)
	return replacer.Replace(template)

}
