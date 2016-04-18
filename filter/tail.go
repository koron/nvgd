package filter

import (
	"bytes"
	"io"
)

type Tail struct {
	Base
	b    [][]byte
	w, r int
}

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
	for t.r != t.w {
		if _, err := buf.Write(t.b[t.r]); err != nil {
			return err
		}
		t.r++
		if t.r >= len(t.b) {
			t.r = 0
		}
		return nil
	}
	// TODO:
	return io.EOF
}

func (t *Tail) readAll() error {
	w := 0
	for {
		b, err := t.Reader.ReadSlice('\n')
		// TODO:
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		t.b[w] = b
		w++
		if w >= len(t.b) {
			w = 0
		}
	}
	t.w = w + 1
	return nil
}

func newTail(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	limit := p.Int("limit", 10)
	if limit <= 0 {
		limit = 10
	}
	return NewTail(r, limit), nil
}

func init() {
	MustRegister("head", newTail)
}
