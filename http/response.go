package http

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

const (
	StatusOK = StatusCode(200)
)

type StatusCode int

func (s StatusCode) String() string {
	switch s {
	case StatusOK:
		return "OK"
	default:
		panic(fmt.Errorf("unsupported status code: %d", s))
	}
}

type Response struct {
	StatusCode StatusCode
	Headers    Headers
	Body       []byte
}

func writeResponse(conn net.Conn, r Response) error {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.Write([]byte("HTTP/1.1"))
	buf.Write(SP)

	buf.Write([]byte(strconv.Itoa(int(r.StatusCode))))
	buf.Write(SP)

	buf.Write([]byte(r.StatusCode.String()))
	buf.Write(CRLF)

	for name, value := range r.Headers {
		buf.Write([]byte(name))
		buf.WriteByte(byte(':'))
		buf.Write(SP)
		buf.Write([]byte(value))
		buf.Write(CRLF)
	}
	buf.Write(CRLF)
	buf.Write(r.Body)

	fmt.Printf("[server]: writing %s to %s\n", buf.Bytes(), conn.RemoteAddr())
	_, err := buf.WriteTo(conn)
	if err != nil {
		return err
	}

	return nil
}

func writeLn(conn net.Conn, bytes []byte) error {
	bytes = append(bytes, CRLF...)
	return write(conn, bytes)
}

func write(conn net.Conn, bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	written := 0
	for written < len(bytes) {
		bytes := bytes[written:]
		fmt.Printf("[server]: writing %b to %s\n", bytes, conn.RemoteAddr())
		n, err := conn.Write(bytes[written:])
		if err != nil {
			return err
		}

		written += n
	}

	return nil
}
