package main

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

// Request is a http Request.
type Request struct {
	Version string
	Method  string
	URL     string
	Headers Headers
}

// Headers represents headers in a http request or response.
type Headers map[string]string

// ReadRequest reads a http request from conn.
func ReadRequest(conn io.Reader) (*Request, error) {
	scanner := bufio.NewScanner(conn)
	var request = Request{Headers: Headers{}}
	var requestLineParsed bool
	var err error
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			break
		}
		if !requestLineParsed {
			request.Method, request.URL, request.Version, err = parseRequestLine(line)
			if err != nil {
				return nil, err
			}
			requestLineParsed = true
			continue
		} else {
			name, value, err := parseHeader(line)
			if err != nil {
				return nil, err
			}
			request.Headers[name] = value
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return &request, nil
}

var requestLineRegex = regexp.MustCompile(`^(.+?) (.+?) HTTP/(.+?)$`)

func parseRequestLine(b []byte) (method, url, version string, err error) {
	m := requestLineRegex.FindSubmatch(b)
	if m == nil {
		err = errors.New("invalid request")
		return
	}
	method = string(m[1])
	url = string(m[2])
	version = string(m[3])
	return
}

var headerRegex = regexp.MustCompile(`^(.+?)\: (.+?)$`)

func parseHeader(b []byte) (name, value string, err error) {
	m := headerRegex.FindSubmatch(b)
	if m == nil {
		err = errors.New("invalid header")
		return
	}
	name, value = string(m[1]), string(m[2])
	return
}
