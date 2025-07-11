package ltsv

import (
	"bufio"
	"bytes"
	"io"
)

// Reader reads LTSV values.
type Reader struct {
	rd *bufio.Reader
}

// NewReader creates a new LTSV reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(r),
	}
}

func NewReaderSize(r io.Reader, maxLineLen int) *Reader {
	return &Reader{
		rd: bufio.NewReaderSize(r, maxLineLen),
	}
}

func (r *Reader) readLine() ([]byte, error) {
	d, err := r.rd.ReadSlice('\n')
	if err == nil || (err == io.EOF && len(d) > 0) {
		return d, nil
	}
	return nil, err
}

// Read read a LTSV value.
func (r *Reader) Read() (*Set, error) {
	d, err := r.readLine()
	if err != nil {
		return nil, err
	}
	d = bytes.TrimLeft(d, " \n\r\t")
	d = bytes.TrimRight(d, "\n\r\t")
	s := &Set{
		Index: make(map[string][]int),
	}
	for _, raw := range bytes.Split(d, []byte("\t")) {
		kv := bytes.SplitN(raw, []byte(":"), 2)
		if len(kv) != 2 {
			continue
		}
		s.Put(string(kv[0]), string(kv[1]))
	}
	return s, nil
}
