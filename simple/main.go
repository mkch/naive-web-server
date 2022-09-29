package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const ListenAddr = ":8080"

func main() {
	log.Println("Starting Server...")
	// Listen for incoming connections.
	ln, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Panicf("Listen failed: %v", err)
	}
	log.Printf("Server started at: %v\n", ListenAddr)
	for {
		// Accept a new connection.
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
		}
		// Handle the new connection.
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection: %v\n", conn.RemoteAddr())

	// Read client request.
	request, err := ReadRequest(conn)
	if err != nil {
		log.Printf("Read request failed: %v\n", err)
		return
	}

	// Handle the request.
	var w = ResponseWriter{
		w:       conn,
		version: request.Version,
		Headers: Headers{"Server": "naive-web-server"}}

	handleRequest(request, &w)
}

func handleRequest(r *Request, w *ResponseWriter) {
	if r.Version != "1.1" {
		w.WriteHeader(505) // HTTP Version Not Supported
		return
	}
	// Backend route map.
	switch r.URL {
	case "/":
		handleIndex(r, w)
	case "/time":
		handleTime(r, w)
	default:
		w.WriteHeader(404)
	}
}

func handleIndex(r *Request, w *ResponseWriter) {
	if r.Method != "GET" {
		w.WriteHeader(405) // Method Not Allowed
		return
	}
	io.WriteString(w, `<html>
<title>This is index</title>
<div>Hello there!</div>
<a href='/time'>Show server time</a>
</html>`)
}

func handleTime(r *Request, w *ResponseWriter) {
	io.WriteString(w, fmt.Sprintf(`<html>
<title>Server time</title>
<span>%v</span>
</html>`, time.Now().Local()))
}
