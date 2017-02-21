package filter

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/resource"
)

type doc struct {
	Headers []string
	Index   map[string]int
	Rows    []row
}

func (d *doc) initHeader(props []ltsv.Property) {
	d.Headers = make([]string, 0, len(props))
	d.Index = make(map[string]int)
	d.Rows = make([]row, 0)
	for _, p := range props {
		d.Index[p.Label] = len(d.Headers)
		d.Headers = append(d.Headers, p.Label)
	}
}

func (d *doc) addRow(props []ltsv.Property) {
	r := row{
		Values: make([]string, len(d.Headers)),
		Others: "",
	}
	for _, p := range props {
		n, ok := d.Index[p.Label]
		if !ok {
			if r.Others != "" {
				r.Others += ", "
			}
			r.Others += p.Label + ":" + p.Value
			continue
		}
		r.Values[n] = p.Value
	}
	if len(r.Others) == 0 {
		r.Others = "(none)"
	}
	d.Rows = append(d.Rows, r)
}

type row struct {
	Values []string
	Others string
}

var tmpl = template.Must(template.New("htmltable").Parse(`<!DOCTYPE html>
<meta charset="UTF-8">
<table border="1">
  <tr>
	{{range .Headers}}
    <th>{{.}}</th>
	{{end}}
	<th>(others)</th>
  </tr>
  {{range .Rows}}
  <tr>
	{{range .Values}}
	<td>{{.}}</td>
	{{end}}
	<td>{{.Others}}</td>
  </tr>
  {{end}}
</table>`))

func filterFunc(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// compose document.
	d := &doc{}
	lr := ltsv.NewReader(r)
	first := true
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
		if first {
			d.initHeader(s.Properties)
			first = false
		}
		d.addRow(s.Properties)
	}
	// execute template.
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, d); err != nil {
		return nil, err
	}
	return r.Wrap(ioutil.NopCloser(buf)), nil
}

func init() {
	filter.MustRegister("htmltable", filterFunc)
}
