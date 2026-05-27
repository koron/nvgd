package core

import (
	"sort"
	"strings"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/resource"
)

// Filters is set of filters.
type Filters struct {
	descs config.FiltersMap
}

func (f *Filters) apply(s *Server, path string, r *resource.Resource) (*resource.Resource, error) {
	filters, found := f.getRaw(path)
	if !found {
		return r, nil
	}
	qp, err := f.parse(filters)
	if err != nil {
		return r, err
	}
	return s.applyFilters(qp, r)
}

func (f *Filters) parse(filters config.Filters) (qparams, error) {
	var qp qparams
	for _, s := range filters {
		p, err := qparamsParse(s)
		if err != nil {
			return qparams{}, err
		}
		if len(p) > 0 {
			qp = append(qp, p...)
		}
	}
	return qp, nil
}

func (f *Filters) getRaw(path string) (config.Filters, bool) {
	// Sort keys by length descending for deterministic longest-prefix matching.
	keys := make([]string, 0, len(f.descs))
	for k := range f.descs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	for _, k := range keys {
		if strings.HasPrefix(path, k) {
			return f.descs[k], true
		}
	}
	return config.Filters{}, false
}
