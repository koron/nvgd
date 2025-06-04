package templateresource

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"path"

	"github.com/koron/nvgd/resource"
)

type TemplateResource struct {
	tmpl *template.Template

	fallback  string
	constData any
}

func New(fs fs.FS, opts ...Option) (*TemplateResource, error) {
	tmpl, err := template.ParseFS(fs, "*/*")
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
	p := u.Path[1:]
	if p == "" {
		p = "index.html"
	}
	bb := &bytes.Buffer{}

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
	for _, path := range candidates {
		tmpl = res.tmpl.Lookup(path)
		if tmpl != nil {
			break
		}
	}
	if tmpl == nil {
		return nil, fmt.Errorf("not found any of: %+s", candidates)
	}

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
