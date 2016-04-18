package filter

import (
	"fmt"
	"io"
	"strconv"
)

type Factory func(io.ReadCloser, Params) (io.ReadCloser, error)

// Filter is abstraction of methods of filtering.
type Filter interface {
	Filter(io.ReadCloser, Params) (io.ReadCloser, error)
}

var filters = map[string]Factory{}

// Register registers a filter with name.
func Register(name string, f Factory) error {
	_, ok := filters[name]
	if ok {
		return fmt.Errorf("duplicated filter name %q", name)
	}
	filters[name] = f
	return nil
}

// MustRegister registers a filter, panic if failed.
func MustRegister(name string, f Factory) {
	if err := Register(name, f); err != nil {
		panic(err)
	}
}

// Find finds a filter.
func Find(name string) Factory {
	f, ok := filters[name]
	if !ok {
		return nil
	}
	return f
}

type Params map[string]string

func (p Params) Int(n string, value int) int {
	s, ok := p[n]
	if !ok {
		return value
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return value
	}
	return v
}
