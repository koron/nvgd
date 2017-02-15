package core

import (
	"io"

	"github.com/koron/nvgd/config"
)

type Filters struct {
	descs config.FiltersMap
}

func (f *Filters) apply(path string, r io.ReadCloser) (io.ReadCloser, error) {
	// TODO:
	return r, nil
}
