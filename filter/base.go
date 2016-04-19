package filter

import (
	"bufio"
	"bytes"
	"io"
)

// Base is base of filters.  It provides common
type Base struct {
	Reader *bufio.Reader
	buf    bytes.Buffer
	closed bool
	raw    io.ReadCloser
	rn     BaseReadNext
}

// BaseReadNext is callback to read next data hunk to buf
type BaseReadNext func(buf *bytes.Buffer) error

// Init initializes Base object.
func (b *Base) Init(r io.ReadCloser, readNext BaseReadNext) {
	b.Reader = bufio.NewReader(r)
	b.raw = r
	b.rn = readNext
}

func (b *Base) Read(buf []byte) (int, error) {
	if b.buf.Len() == 0 {
		if b.closed {
			return 0, io.EOF
		}
		// read next data to b.buf
		b.buf.Reset()
		err := b.rn(&b.buf)
		if err == io.EOF {
			b.Close()
		} else if err != nil {
			return 0, err
		}
	}
	return b.buf.Read(buf)
}

// ReadLine reads a line as []byte.
func (b *Base) ReadLine() ([]byte, error) {
	d, err := b.Reader.ReadSlice('\n')
	if err == nil {
		return d, nil
	} else if err != bufio.ErrBufferFull {
		return nil, err
	}
	bb := bytes.NewBuffer(d)
	for {
		b2, err := b.Reader.ReadSlice('\n')
		if len(b2) > 0 {
			if _, err := bb.Write(b2); err != nil {
				return nil, err
			}
		}
		if err == nil || err == io.EOF {
			return bb.Bytes(), nil
		}
		if err != bufio.ErrBufferFull {
			return nil, err
		}
	}
}

// Close closes head filter.
func (b *Base) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true
	return b.raw.Close()
}
