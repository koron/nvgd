// Package markdown provides markdown filter which render HTML from Markdown.
package markdown

import (
	"bytes"
	"io"
	"regexp"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
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
<meta name="referrer" content="no-referrer">
{{range .Config.CustomCSSURLs}}{{if .}}<link rel="stylesheet" href="{{.}}" type="text/css" />
{{end}}{{end}}
`))

var rxHrefLocalDoc = regexp.MustCompile(`(href="doc/[^."]*\.md)((?:\?[^"]+)?")`)

func appendMarkdownFilter(body []byte) []byte {
	return rxHrefLocalDoc.ReplaceAllFunc(body, func(m []byte) []byte {
		groups := rxHrefLocalDoc.FindSubmatch(m)
		if len(groups) < 3 {
			return m
		}
		beforeQuery := groups[1]
		suffix := groups[2]
		if suffix[0] == '"' {
			// No existing query string.
			return append(append(beforeQuery, []byte("?markdown\"")...), suffix[1:]...)
		}
		// Existing query string: append &markdown before closing quote.
		result := append(append(beforeQuery, suffix[:len(suffix)-1]...), []byte("&markdown\"")...)
		return result
	})
}

func filterMarkdown(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// convert a markdown to HTML as a response body.
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// parse raw as markdown.
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	doc := parser.NewWithExtensions(extensions).Parse(raw)

	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags |
			html.NofollowLinks |
			html.NoreferrerLinks |
			html.NoopenerLinks |
			html.HrefTargetBlank |
			html.FootnoteReturnLinks,
	})
	bodyBytes := markdown.Render(doc, renderer)
	// append "markdown" filter for links to local documents.
	bodyBytes = appendMarkdownFilter(bodyBytes)
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
	return r.Wrap(r2).PutContentType("text/html"), nil
}
