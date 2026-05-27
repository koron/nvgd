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
	"sync"
	"time"
)

type FS struct {
	fs.FS
	mu    sync.RWMutex
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
		return nil, fmt.Errorf("file is directory: %w", fs.ErrNotExist)
	}

	// Check cache for a parsed Template
	tfs.mu.RLock()
	entry, ok := tfs.cache[name]
	tfs.mu.RUnlock()
	if ok && !fi.ModTime().After(entry.ModTime) {
		return entry.Template.Clone()
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
	if _, err := tmpl.Parse(string(body)); err != nil {
		return nil, err
	}

	tfs.mu.Lock()
	tfs.cache[name] = &cacheEntry{
		Template: tmpl,
		ModTime:  fi.ModTime(),
	}
	tfs.mu.Unlock()

	return tmpl.Clone()
}
