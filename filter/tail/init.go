// Package tail provides tail filter.
package tail

import (
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
)

const rtailBufsize = 4096

func newTail(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	limit := p.Int("limit", 10)
	if limit <= 0 {
		limit = 10
	}
	if r2, ok := r.ReadSeekCloser(); ok {
		return r.Wrap(NewRTail(r2, limit, rtailBufsize)), nil
	}
	return r.Wrap(NewTail(r, limit)), nil
}

func init() {
	filter.MustRegister("tail", newTail)
}
