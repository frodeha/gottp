package http

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/frodeha/gottp/buffer"
)

var (
	ErrDoubleClose = errors.New("server is already closed")
)

type Server struct {
	lis    net.Listener
	reader *buffer.Buffer
}

func NewServer(listener net.Listener) *Server {
	return &Server{
		lis: listener,
	}
}

func (s *Server) ServeHTTP() error {
	fmt.Printf("[server]: serving HTTP on %s\n", s.lis.Addr())
	for {
		conn, err := s.lis.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		s.handle(conn)
	}
}

const (
	MethodOptions = "OPTIONS"
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodTrace   = "TRACE"
	MethodConnect = "CONNECT"

	HTTP1_1 = "HTTP/1.1"
)

func (s *Server) loanBufferReader() *buffer.Buffer {
	if s.reader == nil {
		s.reader = buffer.New()
	}
	return s.reader
}

func (s *Server) handle(conn net.Conn) {
	fmt.Printf("[server]: handling connection from %s\n", conn.RemoteAddr())
	defer conn.Close()

	reader := s.loanBufferReader()
	reader.SetReader(conn)

	req, err := parseRequest(reader)
	if err != nil {
		fmt.Printf("[server]: error while parsing request from %s: %s\n", conn.RemoteAddr(), err)
		return
	}
	fmt.Printf("[server]: parsed request: %+v\n", req)

	res := Response{
		StatusCode: StatusOK,
		Headers: Headers{
			"Content-Type":   "text/plain",
			"Content-Length": strconv.Itoa(len(req.Body)),
		},
		Body: req.Body,
	}

	err = writeResponse(conn, res)
	if err != nil {
		fmt.Printf("[server]: error while writing response to %s: %s\n", conn.RemoteAddr(), err)
		return
	}
}

var (
	CRLF = []byte{0x0D, 0x0A}
	SP   = []byte{0x20}
)

func (s *Server) Close() error {
	fmt.Printf("[server]: closing\n")
	return s.lis.Close()
}
