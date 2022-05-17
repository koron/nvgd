// Package indexhtml provides index HTML filter.
package indexhtml

import (
	"bytes"
	"html/template"
	"io"
	"path"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/internal/ltsv"
	"github.com/koron/nvgd/resource"
)

var tmpl = template.Must(template.New("indexhtml").Parse(`<!DOCTYPE! html>
<meta charset="UTF-8">
{{range .Config.CustomCSSURLs}}{{if .}}<link rel="stylesheet" href="{{.}}" type="text/css" />
{{end}}{{end}}
<div>
  {{if .UpLink}}<a href="{{.UpLink}}">Up</a>{{end}}
  {{if .NextLink}}<a href="{{.NextLink}}">Next</a>{{end}}
</div>
<table border="1">
  <tr><th>Name</th><th>Type</th><th>Size</th><th>Modified At</th><th>Download</th></tr>
  {{range .Entries}}
  <tr>
    <td><a href="{{.Link}}">{{.Name}}</a></td>
    <td>{{.Type}}</td>
    <td>{{.Size}}</td>
    <td>{{.ModifiedAt}}</td>
	<td>{{if .Download}}<a href="{{.Download}}">DL</a>{{end}}</td>
  </tr>
  {{end}}
</table>`))

type doc struct {
	Entries  []entry
	UpLink   string
	NextLink string

	Config *Config
}

type entry struct {
	Name       string
	Type       string
	Size       string
	ModifiedAt string
	Link       string
	Download   string
}

type Config struct {
	CustomCSSURLs []string `yaml:"custom_css_urls,omitempty"`
}

var cfg Config

func pathPrefix(s string) string {
	if s == "" {
		return ""
	}
	return path.Join(config.Root().PathPrefix, s)
}

func filterFunc(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// compose document.
	d := &doc{
		Config: &cfg,
	}
	lr := ltsv.NewReader(r)
	for {
		s, err := lr.Read()
		if err != nil {
			r.Close()
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if s.Empty() {
			continue
		}
		e := entry{
			Name:       s.GetFirst("name"),
			Type:       s.GetFirst("type"),
			Size:       s.GetFirst("size"),
			ModifiedAt: s.GetFirst("modified_at"),
			Link:       pathPrefix(s.GetFirst("link")),
			Download:   pathPrefix(s.GetFirst("download")),
		}
		d.Entries = append(d.Entries, e)
	}
	if link, ok := r.String(commonconst.UpLink); ok {
		d.UpLink = pathPrefix(link)
	}
	if link, ok := r.String(commonconst.NextLink); ok {
		d.NextLink = pathPrefix(link)
	}
	// execute template.
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, d); err != nil {
		return nil, err
	}
	return r.Wrap(io.NopCloser(buf)), nil
}

func init() {
	filter.MustRegister("indexhtml", filterFunc)
	config.RegisterFilter("indexhtml", &cfg)
}
