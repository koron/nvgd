// Package resource provides data transfer objects for NVGD's protocol and
// filters.
package resource

import (
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/koron/nvgd/internal/commonconst"
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

// NewRedirect creates a Resource to redirect to another path.
func NewRedirect(redirectPath string) *Resource {
	r := NewString("redirect to: " + redirectPath)
	r.Options[commonconst.Redirect] = redirectPath

	return r
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
func (r *Resource) Put(name string, value any) *Resource {
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

// PutFilename puts a filename option.
func (r *Resource) PutFilename(s string) *Resource {
	return r.PutString(Filename, path.Base(s)).GuessContentType(s)
}

func (r *Resource) PutAttachmentFilename(s string) *Resource {
	return r.PutString(AttachmentFilename, s).GuessContentType(s)
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
