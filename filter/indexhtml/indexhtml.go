// Package indexhtml provides index HTML filter.
package indexhtml

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

var tmpl = template.Must(template.New("indexhtml").Parse(`<!DOCTYPE html>
<meta charset="UTF-8">
{{range .Config.CustomCSSURLs}}{{if .}}<link rel="stylesheet" href="{{.}}" type="text/css" />
{{end}}{{end}}
<div>
  {{if .UpLink}}<a href="{{.UpLink}}">Up</a>{{end}}
  {{if .NextLink}}<a href="{{.NextLink}}">Next</a>{{end}}
</div>
<table border="1">
  <tr><th>Name</th><th>Type</th><th>Size</th><th>Modified At</th><th>Actions</th></tr>
  {{range .Entries}}
  <tr>
    <td><a href="{{.Link}}">{{.Name}}</a></td>
    <td>{{.Type}}</td>
    <td>{{.Size}}</td>
    <td>{{.ModifiedAt}}</td>
	<td>
		{{if .Download}}<a href="{{.Download}}">DL</a>{{end}}
		{{if .QueryLink}}<a href="{{.QueryLink}}">Query</a>{{end}}
		{{if .DuckDBLink}}<a href="{{.DuckDBLink}}">DuckDB</a>{{end}}
	</td>
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
	QueryLink  string
	DuckDBLink string
}

type Config struct {
	CustomCSSURLs []string `yaml:"custom_css_urls,omitempty"`
}

var cfg Config

func pathPrefix(s string) string {
	if s == "" {
		return ""
	}
	// path.Join cleans "//" at "/". it break some links. so we can't use it.
	//return path.Join(config.Root().PathPrefix, s)
	return strings.TrimRight(config.Root().PathPrefix, "/") + "/" + strings.TrimLeft(s, "/")
}

func chooseTimeLayout(name string) string {
	switch strings.ToUpper(name) {
	case "ANSIC":
		return time.ANSIC
	case "UNIX":
		return time.UnixDate
	case "RUBY":
		return time.RubyDate
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC850":
		return time.RFC850
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "RFC3339":
		return time.RFC3339
	case "RFC3339NANO":
		return time.RFC3339Nano
	case "STAMP":
		return time.Stamp
	case "DATETIME":
		return time.DateTime
	default:
		return time.RFC1123
	}
}

func filterFunc(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	timeLayout := chooseTimeLayout(p.String("timefmt", "RFC1123"))
	noUpLink := p.Bool("nouplink", false)
	// compose document.
	d := &doc{
		Config: &cfg,
	}
	lr := filterbase.NewLTSVReader(r)
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
		name := s.GetFirst("name")
		e := entry{
			Name:       name,
			Type:       s.GetFirst("type"),
			Size:       s.GetFirst("size"),
			ModifiedAt: s.GetFirst("modified_at"),
			Link:       pathPrefix(s.GetFirst("link")),
			Download:   pathPrefix(s.GetFirst("download")),
		}
		if fmt, ok := supportQuery(name); ok {
			qlink := "/trdsql/"
			qlink += "s=" + url.PathEscape(pathPrefix(s.GetFirst("link")))
			qlink += "&q=" + url.PathEscape("SELECT * FROM t")
			qlink += "&ifmt=" + fmt
			qlink += "&ih=false"
			e.QueryLink = qlink
		}
		if _, ok := supportDuckDB(name); ok {
			link := "/duckdb/show-as-view?"
			link += "t=" + url.PathEscape(pathPrefix(s.GetFirst("link")))
			e.DuckDBLink = link
		}
		// detect UNIX time, and convert it to specified time layout (default is RFC1123).
		if sec, err := strconv.ParseInt(e.ModifiedAt, 10, 64); err == nil {
			e.ModifiedAt = time.Unix(sec, 0).Format(timeLayout)
		}
		d.Entries = append(d.Entries, e)
	}
	if link, ok := r.String(commonconst.UpLink); ok && !noUpLink {
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
	return r.Wrap(io.NopCloser(buf)).PutContentType("text/html"), nil
}

// supportQuery checks filename, which is supported format as query target.
func supportQuery(name string) (string, bool) {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".csv", ".ltsv", ".tsv":
		return strings.ToUpper(ext[1:]), true
	default:
		return "", false
	}
}

// supportDuckDB checks filename, which is supported by DuckDB.
func supportDuckDB(name string) (string, bool) {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".csv", ".xlsx", ".json", ".parquet":
		// supported file formats came from:
		// https://duckdb.org/docs/stable/guides/file_formats/overview
		return strings.ToUpper(ext[1:]), true
	default:
		return "", false
	}
}

func init() {
	filter.MustRegister("indexhtml", filterFunc)
	config.RegisterFilter("indexhtml", &cfg)
}
