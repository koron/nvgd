package filter

import (
	"bytes"
	"io"
)

func init() {
	MustRegister("head", &headFactory{})
}

type headFactory struct {
}

func (f *headFactory) Filter(r io.ReadCloser, p map[string]string) (io.ReadCloser, error) {
	// TODO:
	return NewHead(r, 0, 10), nil
}

// Head is "head" like filter.
type Head struct {
	buf    bytes.Buffer
	closed bool
	reader io.ReadCloser
	start  uint
	limit  uint
}

var (
	_ io.ReadCloser = (*Head)(nil)
)

// NewHead creates an instance of head filter.
func NewHead(r io.ReadCloser, start, limit uint) *Head {
	return &Head{
		reader: r,
		start:  start,
		limit:  limit,
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
	// TODO:
	return nil
}

// Close closes head filter.
func (h *Head) Close() error {
	if h.closed {
		return nil
	}
	h.closed = true
	return h.reader.Close()
}
