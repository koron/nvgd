package filter

import (
	"fmt"
	"strconv"

	"github.com/koron/nvgd/resource"
)

// Factory is filter factory.
type Factory func(*resource.Resource, Params) (*resource.Resource, error)

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

// Params represents parameters for filter.
type Params map[string]string

// Int gets int value from Params by name.
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

// String gets string value from Params by name.
func (p Params) String(n string, value string) string {
	s, ok := p[n]
	if !ok {
		return value
	}
	return s
}

// Bool gets bool value from Params by name.
func (p Params) Bool(n string, value bool) bool {
	s, ok := p[n]
	if !ok {
		return value
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return value
	}
	return v
}
