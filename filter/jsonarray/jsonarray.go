// Package jsonarray provides JSON array filter.
package jsonarray

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

func init() {
	filter.MustRegister("jsonarray", newFilter)
}

func newFilter(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	return r.Wrap(New(r)), nil
}

type Filter struct {
	filterbase.Base

	first bool
	last  bool
}

func New(r io.ReadCloser) *Filter {
	f := &Filter{}
	f.Base.Init(r, f.readNext)
	return f
}

func (f *Filter) readNext(buf *bytes.Buffer) error {
	if f.last {
		return io.EOF
	}
	if !f.first {
		buf.WriteByte('[')
		f.first = true
	}
	b, err0 := f.Base.ReadLine()
	if n := len(b); n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	j, err := json.Marshal(string(b))
	if err != nil {
		return err
	}
	buf.Write(j)
	switch err0 {
	case nil:
		buf.WriteString(",\n")
	case io.EOF:
		buf.WriteString("]\n")
		f.last = true
	}
	return err0
}
