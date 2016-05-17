package lz4

import (
	"errors"
	"io"
)

// Reader reades lz4 decoded stream.
type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) (*Reader, error) {
	return &Reader{
		r: r,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	// TODO:
	return 0, errors.New("lz4.Reader.Read() is not implemented yet")
}

func (r *Reader) Close() error {
	if c, ok := r.r.(io.ReadCloser); ok {
		return c.Close()
	}
	return nil
}
