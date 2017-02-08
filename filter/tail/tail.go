package tail

import (
	"bytes"
	"io"

	"github.com/koron/nvgd/filter"
)

// Tail is "tail" like filter.
type Tail struct {
	filter.Base
	b    [][]byte
	w, r int
}

// NewTail creates an instance of tail filter.
func NewTail(r io.ReadCloser, limit int) *Tail {
	t := &Tail{
		b: make([][]byte, limit),
		w: -1,
	}
	t.Base.Init(r, t.readNext)
	return t
}

func (t *Tail) readNext(buf *bytes.Buffer) error {
	if t.w < 0 {
		if err := t.readAll(); err != nil {
			return err
		}
	}
	if t.r == t.w {
		return io.EOF
	}
	for {
		if b := t.b[t.r]; b != nil {
			if _, err := buf.Write(b); err != nil {
				return err
			}
		}
		if t.r == t.w {
			return io.EOF
		}
		t.r = t.addr(t.r + 1)
	}
}

func (t *Tail) readAll() error {
	for {
		b, err := t.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		t.putLine(b)
	}
	// setup read pointer
	t.r = t.addr(t.w + 1)
	return nil
}

func (t *Tail) putLine(b []byte) {
	t.w = t.addr(t.w + 1)
	t.b[t.w] = b
}

func (t *Tail) addr(n int) int {
	l := len(t.b)
	for n >= l {
		n -= l
	}
	return n
}
