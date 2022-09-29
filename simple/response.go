package main

import (
	"fmt"
	"io"
	"log"
)

// A ResponseWriter is used by an HTTP handler to construct an HTTP response.
type ResponseWriter struct {
	Headers Headers

	version       string
	w             io.Writer
	headerWritten bool
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.headerWritten {
		log.Printf("Header already written\n")
		return
	}
	w.headerWritten = true
	var status string
	switch statusCode {
	case 200:
		status = "OK"
	case 404:
		status = "Not Found"
	case 405:
		status = "Method Not Allowed"
	case 505:
		status = "HTTP Version Not Supported"
	default:
		log.Printf("Invalid status code: %v\n", statusCode)
	}

	// Write response line.
	_, err := fmt.Fprintf(w.w, "HTTP/%v %v %v\r\n", w.version, statusCode, status)
	if err != nil {
		log.Printf("Write response failed: %v\n", err)
	}

	// Write response headers.
	for k, v := range w.Headers {
		_, err = fmt.Fprintf(w.w, "%v: %v\r\n", k, v)
		if err != nil {
			log.Printf("Write response failed: %v\n", err)
		}
	}

	// Write "EOF"
	_, err = w.w.Write([]byte("\r\n"))
	if err != nil {
		log.Printf("Write response failed: %v\n", err)
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.headerWritten {
		w.WriteHeader(200)
	}
	return w.w.Write(b)
}
