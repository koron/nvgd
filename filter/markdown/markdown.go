// Package markdown provides markdown filter which render HTML from Markdown.
package markdown

import (
	"bytes"
	"io"
	"text/template"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
	"github.com/russross/blackfriday/v2"
)

type Config struct {
	CustomCSSURLs []string `yaml:"custom_css_urls,omitempty"`
}

var cfg Config

func init() {
	filter.MustRegister("markdown", filterMarkdown)
	config.RegisterFilter("markdown", &cfg)
}

var tmpl = template.Must(template.New("markdown").Parse(`<!DOCTYPE html>
<meta charset="UTF-8">
{{range .Config.CustomCSSURLs}}{{if .}}<link rel="stylesheet" href="{{.}}" type="text/css" />
{{end}}{{end}}
`))

func filterMarkdown(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// convert a markdown to HTML as a response body.
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags |
			blackfriday.NofollowLinks |
			blackfriday.NoreferrerLinks |
			blackfriday.NoopenerLinks |
			blackfriday.HrefTargetBlank |
			blackfriday.FootnoteReturnLinks,
	})
	extensions := blackfriday.CommonExtensions |
		blackfriday.AutoHeadingIDs
	bodyBytes := blackfriday.Run(raw,
		blackfriday.WithExtensions(extensions),
		blackfriday.WithRenderer(renderer))
	// generate header
	d := struct {
		Config *Config
	}{
		Config: &cfg,
	}
	head := new(bytes.Buffer)
	if err := tmpl.Execute(head, d); err != nil {
		return nil, err
	}
	r2 := io.NopCloser(io.MultiReader(head, bytes.NewReader(bodyBytes)))
	return r.Wrap(r2), nil
}
