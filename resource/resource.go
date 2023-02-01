// Package resource provides data transfer objects for NVGD's protocol and
// filters.
package resource

import (
	"bytes"
	"io"
	"path"
	"strings"
)

// Resource packs ReadCloser and its meta info.
type Resource struct {
	io.ReadCloser
	Options
}

// New creates a Resource from ReadCloser.
func New(rc io.ReadCloser) *Resource {
	return &Resource{
		ReadCloser: rc,
		Options:    Options{},
	}
}

// NewString creates a Resource with string as content.
func NewString(s string) *Resource {
	b := bytes.NewReader([]byte(s))
	return New(io.NopCloser(b))
}

// Raw returns underlying io.ReadCloser in this resource.
func (r *Resource) Raw() io.ReadCloser {
	return r.ReadCloser
}

// Wrap creates new Resource with a io.ReadCloser which inherits properties
// from current Resource.
func (r *Resource) Wrap(rc io.ReadCloser) *Resource {
	// TODO: copy other properties of r.
	return &Resource{
		ReadCloser: rc,
		Options:    r.Options.clone(),
	}
}

// ReadSeekCloser obtains ReadSeekCloser if it could.
func (r *Resource) ReadSeekCloser() (ReadSeekCloser, bool) {
	x, ok := r.ReadCloser.(ReadSeekCloser)
	return x, ok
}

// Put puts a pair of name and value as an option.
func (r *Resource) Put(name string, value interface{}) *Resource {
	r.Options[name] = value
	return r
}

// PutString puts a string as an option. When value is empty string, it deletes
// the option.
func (r *Resource) PutString(name, value string) *Resource {
	if value == "" {
		delete(r.Options, name)
	} else {
		r.Options[name] = value
	}
	return r
}

// PutFilename puts a filenaem option.
func (r *Resource) PutFilename(s string) *Resource {
	return r.PutString(Filename, s).GuessContentType(s)
}

// PutContentType puts a content-type option.
func (r *Resource) PutContentType(s string) *Resource {
	return r.PutString(ContentType, s)
}

// GuessContentType guess a content-type from argument string.
func (r *Resource) GuessContentType(s string) *Resource {
	ex := strings.ToLower(path.Ext(s))
	if ct, ok := Mime[ex]; ok {
		r.PutContentType(ct)
	}
	return r
}
