// Package templatefs provides Template set.
//
// This template set is linked to fs.FS and returns a template parsed from the
// specified file in fs.FS.
package templatefs

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"time"
)

type FS struct {
	fs.FS
	cache map[string]*cacheEntry
}

type cacheEntry struct {
	*template.Template
	ModTime time.Time
}

func New(fsys fs.FS) *FS {
	return &FS{
		FS:    fsys,
		cache: map[string]*cacheEntry{},
	}
}

type Option interface {
	apply(*template.Template) (*template.Template, error)
}

type OptionFunc func(*template.Template) (*template.Template, error)

func (fn OptionFunc) apply(tmpl *template.Template) (*template.Template, error) {
	return fn(tmpl)
}

var _ Option = (OptionFunc)(nil)

type options []Option

func (opts options) apply(tmpl *template.Template) (*template.Template, error) {
	var err error
	for _, opt := range opts {
		tmpl, err = opt.apply(tmpl)
		if err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}

var _ Option = (options)(nil)

func (tfs *FS) Template(name string, opts ...Option) (*template.Template, error) {
	f, err := tfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("file is direcotry: %w", fs.ErrNotExist)
	}

	// Check cache for a parsed Template
	if entry, ok := tfs.cache[name]; ok && !fi.ModTime().After(entry.ModTime) {
		return entry.Template, nil
	}

	// Create a new template and apply options
	tmpl, err := options(opts).apply(template.New(name))
	if err != nil {
		return nil, err
	}

	// Parse body of the template
	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tmpl.Parse(string(body))

	tfs.cache[name] = &cacheEntry{
		Template: tmpl,
		ModTime:  fi.ModTime(),
	}
	return tmpl, nil
}
