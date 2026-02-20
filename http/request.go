package http

import (
	"fmt"
	"strings"

	"github.com/frodeha/gottp/buffer"
)

type Request struct {
	Protocol string
	Method   string
	URL      string
	Headers  Headers
	Body     []byte
}

type Headers map[string]string

func parseRequest(r *buffer.BytesBuffer) (Request, error) {
	var (
		req Request
	)

	method, err := r.ReadString(SP)
	if err != nil {
		return req, err
	}

	switch method {
	case MethodOptions, MethodGet, MethodHead, MethodPost, MethodPut, MethodDelete, MethodTrace, MethodConnect:
		req.Method = method
	default:
		return req, fmt.Errorf("invalid method: %s", method)
	}

	url, err := r.ReadString(SP)
	if err != nil {
		return req, err
	}
	req.URL = url

	protocol, err := r.ReadString(CRLF)
	if err != nil {
		return req, err
	}

	if string(protocol) != HTTP1_1 {
		return req, fmt.Errorf("unsupported protocol: %s", protocol)
	}
	req.Protocol = protocol

	req.Headers = make(Headers)
	for {
		line, err := r.ReadString(CRLF)
		if err != nil {
			return req, err
		}

		if line == "" {
			break // We're done
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return req, fmt.Errorf("invalid header line: %s", line)
		}

		name := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		req.Headers[name] = value
	}

	return req, nil
}
