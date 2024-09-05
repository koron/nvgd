package filter

import (
	"bufio"
	"bytes"
	"io"

	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

func init() {
	MustRegister("head", newHead)
}

func newHead(r *resource.Resource, p Params) (*resource.Resource, error) {
	start := p.Int("start", 0)
	if start < 0 {
		start = 0
	}
	limit := p.Int("limit", 10)
	if limit <= 0 {
		limit = 10
	}
	return r.Wrap(NewHead(r, uint(start), uint(limit))), nil
}

// Head is "head" like filter.
type Head struct {
	filterbase.Base
	start uint
	last  uint
	curr  uint
}

var (
	_ io.ReadCloser = (*Head)(nil)
)

// NewHead creates an instance of head filter.
func NewHead(r io.ReadCloser, start, limit uint) *Head {
	h := &Head{
		start: start,
		last:  start + limit,
	}
	h.Base.Init(r, h.readNext)
	return h
}

func (h *Head) readNext(buf *bytes.Buffer) error {
	for h.curr < h.last {
		lnum := h.curr
		b, err := h.Reader.ReadSlice('\n')
		if err != nil {
			if err != bufio.ErrBufferFull {
				return err
			}
		} else {
			h.curr++
		}
		if lnum >= h.start {
			if _, err := buf.Write(b); err != nil {
				return err
			}
			return nil
		}
	}
	return io.EOF
}
