package count

import (
	"bytes"
	"io"
	"strconv"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
)

func init() {
	filter.MustRegister("count", newCount)
}

func newCount(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	return r.Wrap(New(r)), nil
}

type Count struct {
	filter.Base
	n int64
}

func New(r io.ReadCloser) *Count {
	c := &Count{n: -1}
	c.Base.Init(r, c.readNext)
	return c
}

func (c *Count) readNext(buf *bytes.Buffer) error {
	if c.n >= 0 {
		return io.EOF
	}
	c.n = 0
	for {
		raw, err := c.ReadLine()
		if err != nil && len(raw) == 0 {
			if err == io.EOF {
				if _, err := buf.WriteString(strconv.FormatInt(c.n, 10)); err != nil {
					return err
				}
				return nil
			}
			return err
		}
		c.n++
	}
}
