package templateresource

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"path"

	"github.com/koron/nvgd/internal/templatefs"
	"github.com/koron/nvgd/resource"
)

type TemplateResource struct {
	fs *templatefs.FS

	fallback  string
	constData any
}

func New(fsys fs.FS, opts ...Option) (*TemplateResource, error) {
	r := &TemplateResource{
		fs:       templatefs.New(fsys),
		fallback: "index.html",
	}
	for _, opt := range opts {
		opt.apply(r)
	}
	return r, nil
}

func (res *TemplateResource) Open(u *url.URL) (*resource.Resource, error) {
	if u.Path == "" {
		u.Path = "/"
		return resource.NewRedirect(u.String()), nil
	}
	p := u.Path[1:]
	if p == "" {
		p = "index.html"
	}

	// compose data for template
	data := map[string]any{}
	if res.constData != nil {
		data["constant"] = res.constData
	}
	data["query"] = u.Query()

	// Lookup a template.
	var (
		candidates []string
		tmpl       *template.Template
	)
	candidates = append(candidates, p)
	if path.Ext(p) == "" {
		candidates = append(candidates, p+".html")
	}
	if res.fallback != "" {
		candidates = append(candidates, res.fallback)
	}
	for _, name := range candidates {
		var err error
		tmpl, err = res.fs.Template(name)
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	if tmpl == nil {
		return nil, fmt.Errorf("not found any of: %+s", candidates)
	}

	bb := &bytes.Buffer{}
	err := tmpl.Execute(bb, data)
	if err != nil {
		return nil, err
	}
	return resource.New(io.NopCloser(bb)).GuessContentType(p), nil
}

type Option interface {
	apply(*TemplateResource)
}

type optionFunc func(*TemplateResource)

func (fn optionFunc) apply(res *TemplateResource) {
	fn(res)
}

func WithConstant(data any) Option {
	return optionFunc(func(res *TemplateResource) {
		res.constData = data
	})
}
