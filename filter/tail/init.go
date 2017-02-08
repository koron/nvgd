package tail

import (
	"io"

	"github.com/koron/nvgd/filter"
)

const rtailBufsize = 4096

func newTail(r io.ReadCloser, p filter.Params) (io.ReadCloser, error) {
	limit := p.Int("limit", 10)
	if limit <= 0 {
		limit = 10
	}
	if r2, ok := r.(readSeekCloser); ok {
		return NewRTail(r2, limit, rtailBufsize), nil
	}
	return NewTail(r, limit), nil
}

func init() {
	filter.MustRegister("tail", newTail)
}
