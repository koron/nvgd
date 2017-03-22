package resource

import (
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

func (r *Resource) Put(name string, value interface{}) *Resource {
	r.Options[name] = value
	return r
}

func (r *Resource) PutContentType(s string) *Resource {
	if s == "" {
		delete(r.Options, ContentType)
	} else {
		r.Options[ContentType] = s
	}
	return r
}

func (r *Resource) GuessContentType(s string) *Resource {
	ex := strings.ToLower(path.Ext(s))
	if ct, ok := Mime[ex]; ok {
		r.PutContentType(ct)
	}
	return r
}
