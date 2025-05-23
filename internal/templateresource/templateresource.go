package templateresource

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"

	"github.com/koron/nvgd/resource"
)

type TemplateResource struct {
	tmpl *template.Template

	fallback  string
	constData any
}

func New(fs fs.FS, opts ...Option) (*TemplateResource, error) {
	tmpl, err := template.ParseFS(fs, "**/*")
	if err != nil {
		return nil, err
	}
	r := &TemplateResource{
		tmpl:     tmpl,
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
	path := u.Path
	if path == "/" {
		path = "index.html"
	}
	bb := &bytes.Buffer{}

	// compose data for template
	data := map[string]any{}
	if res.constData != nil {
		data["constant"] = res.constData
	}

	tmpl := res.tmpl.Lookup(path)
	if tmpl == nil {
		if res.fallback == "" {
			return nil, fmt.Errorf("not found: %s", path)
		}
		tmpl = res.tmpl.Lookup(res.fallback)
		if tmpl == nil {
			return nil, fmt.Errorf("not found: %s and %s", path, res.fallback)
		}
	}
	err := res.tmpl.Execute(bb, data)
	if err != nil {
		return nil, err
	}
	return resource.New(io.NopCloser(bb)).GuessContentType(path), nil
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
