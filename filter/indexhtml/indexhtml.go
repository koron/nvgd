package indexhtml

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/ltsv"
)

var tmpl = template.Must(template.New("indexhtml").Parse(`<!DOCTYPE! html>
<meta charset="UTF-8">
<table border="1">
  <tr><th>Name</th><th>Type</th><th>Size</th><th>Modified At</th></tr>
  {{range .Entries}}
  <tr>
    <td><a href="{{.Link}}">{{.Name}}</a></td>
    <td>{{.Type}}</td>
    <td>{{.Size}}</td>
    <td>{{.ModifiedAt}}</td>
  </tr>
  {{end}}
</table>`))

type doc struct {
	Entries []entry
}

type entry struct {
	Name       string
	Type       string
	Size       string
	ModifiedAt string
	Link       string
}

func filterFunc(r io.ReadCloser, p filter.Params) (io.ReadCloser, error) {
	// compose document.
	d := &doc{}
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
			Link:       s.GetFirst("link"),
		}
		d.Entries = append(d.Entries, e)
	}
	// execute template.
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, d); err != nil {
		return nil, err
	}
	return ioutil.NopCloser(buf), nil
}

func init() {
	filter.MustRegister("indexhtml", filterFunc)
}
