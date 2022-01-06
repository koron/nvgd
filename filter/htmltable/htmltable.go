package filter

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"time"

	"github.com/koron/nvgd/common_const"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/resource"
)

type doc struct {
	Headers []string
	Index   map[string]int
	Rows    []row

	HasOthers  bool
	HasOptions bool

	SQLQuery       *string
	SQLTruncatedBy *int
	SQLExecTime    *string
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
	if r.Others == "" {
		r.Others = "(none)"
	} else {
		d.HasOthers = true
	}
	d.Rows = append(d.Rows, r)
}

type row struct {
	Values []string
	Others string
}

var tmpl = template.Must(template.New("htmltable").Parse(`<!DOCTYPE html>
<meta charset="UTF-8">
<style>
td, th {
  white-space: nowrap;
  padding: 2px 10px;
  text-align: left;
}
table {
  border-collapse: collapse;
}
textarea#query {
  width: 100%;
  max-width: 100%;
}
</style>
{{if .HasOptions}}
<dl>
{{if .SQLQuery}}<dt>Statement (SQL)</dt><dd><textarea id="query" readonly rows="12">{{.SQLQuery}}</textarea><br><button id="edit">Edit</button></dd>{{end}}
{{if .SQLExecTime}}<dt>Execution time</dt><dd><code>{{.SQLExecTime}}</code></dd>{{end}}
{{if .SQLTruncatedBy}}<dt><code>max_rows</code> applied (SQL)</dt><dd>only <code>{{.SQLTruncatedBy}}</code> rows are shown</dd>{{end}}
</dl>
{{end}}
<table border="1">
  <tr>
	{{range .Headers}}
    <th>{{.}}</th>
	{{end}}
	{{if .HasOthers}}
	<th>(others)</th>
	{{end}}
  </tr>
  {{range .Rows}}
  <tr>
	{{range .Values}}
	<td>{{.}}</td>
	{{end}}
	{{if $.HasOthers}}
	<td>{{.Others}}</td>
	{{end}}
  </tr>
  {{end}}
</table>
<script>
(function(g) {
  'use strict'
  var d = g.document;
  var query = d.querySelector('#query');
  var edit = d.querySelector('#edit');
  edit.addEventListener('click', function(ev) {
	ev.preventDefault();
	g.sessionStorage.setItem('query', query.value);
	var url = g.location.href;
	g.location.href = url.slice(0, url.lastIndexOf('/')+1);
  });
})(this);
</script>
`))

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
	if v, ok := r.String(common_const.SQLQuery); ok {
		d.SQLQuery = &v
		d.HasOptions = true
	}
	if v, ok := r.Int(common_const.SQLTruncatedBy); ok {
		d.SQLTruncatedBy = &v
		d.HasOptions = true
	}
	if v, ok := r.Options[common_const.SQLExecTime]; ok {
		if w, ok := v.(time.Duration); ok {
			s := w.String()
			d.SQLExecTime = &s
			d.HasOptions = true
		}
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
