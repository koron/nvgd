package filterbase

import (
	"bufio"
	"errors"
	"io"
)

type LineReader struct {
	buf *bufio.Reader
}

func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{
		buf: bufio.NewReaderSize(r, Config.MaxLineLen),
	}
}

func (r *LineReader) ReadLine() ([]byte, error) {
	b, err := r.buf.ReadSlice('\n')
	if err == nil || errors.Is(err, io.EOF) {
		bb := make([]byte, len(b))
		copy(bb, b)
		return bb, err
	}
	if errors.Is(err, bufio.ErrBufferFull) {
		return nil, ErrMaxLineExceeded
	}
	return nil, err
}
