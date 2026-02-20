package buffer

import (
	"errors"
	"fmt"
	"io"
)

type Buffer struct {
	r   io.Reader
	buf []byte

	bytes  int
	cursor int
}

func (b *Buffer) ReadString(delim []byte) (string, error) {
	if b.r == nil {
		return "", errors.New("missing reader")
	}

	if b.bytes == 0 {
		err := b.fill()
		if err != nil {
			return "", err
		}
	}

	to := b.cursor
	for {
		for to < b.bytes {
			if b.buf[to] != delim[0] {
				to++
				continue
			}

			matches := true
			for i := 1; i < len(delim) && to+i < len(b.buf); i++ {
				matches = matches && b.buf[to+i] == delim[i]
			}

			if matches {
				out := b.buf[b.cursor:to]
				b.cursor = to + len(delim)
				return string(out), nil
			}
		}

		if b.bytes == len(b.buf) {
			b.expand()
			err := b.fill()
			if err != nil {
				return "", err
			}
			continue
		}

		return "", io.EOF
	}
}

func (b *Buffer) SetReader(r io.Reader) {
	b.r = r
	b.cursor = 0
	b.bytes = 0
}

func (b *Buffer) expand() {
	new := make([]byte, len(b.buf)*2)
	copy(new, b.buf)
	b.buf = new
}

func (b *Buffer) fill() error {
	c, err := b.r.Read(b.buf[b.bytes:])
	if err != nil {
		return err
	}
	b.bytes = b.bytes + c

	fmt.Printf("[server]: read %d bytes \"%s\"\n", c, b.buf)
	return nil
}

func New() *Buffer {
	return &Buffer{
		buf:    make([]byte, 128),
		cursor: 0,
		bytes:  0,
	}
}
