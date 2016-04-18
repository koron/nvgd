package filter

import (
	"bufio"
	"bytes"
	"io"
)

func init() {
	MustRegister("head", newHead)
}

func newHead(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	start := p.Int("start", 0)
	if start < 0 {
		start = 0
	}
	limit := p.Int("limit", 10)
	if limit <= 0 {
		limit = 10
	}
	return NewHead(r, uint(start), uint(limit)), nil
}

// Head is "head" like filter.
type Head struct {
	buf    bytes.Buffer
	closed bool
	raw    io.ReadCloser
	rd     *bufio.Reader
	start  uint
	last   uint
	curr   uint
}

var (
	_ io.ReadCloser = (*Head)(nil)
)

// NewHead creates an instance of head filter.
func NewHead(r io.ReadCloser, start, limit uint) *Head {
	return &Head{
		raw:   r,
		rd:    bufio.NewReader(r),
		start: start,
		last:  start + limit,
	}
}

func (h *Head) Read(buf []byte) (int, error) {
	if h.buf.Len() == 0 {
		if h.closed {
			return 0, io.EOF
		}
		// read next data to h.buf
		h.buf.Reset()
		err := h.readNext(&h.buf)
		if err == io.EOF {
			h.Close()
		} else if err != nil {
			return 0, err
		}
	}
	return h.buf.Read(buf)
}

func (h *Head) readNext(buf *bytes.Buffer) error {
	for h.curr < h.last {
		lnum := h.curr
		b, err := h.rd.ReadSlice('\n')
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

// Close closes head filter.
func (h *Head) Close() error {
	if h.closed {
		return nil
	}
	h.closed = true
	return h.raw.Close()
}
