package filter

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/koron/nvgd/ltsv"
)

var htTmpl1 = template.Must(template.New("htmltable").Parse(`<!DOCTYPE! html>
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

type htDoc struct {
	Entries []htEntry
}

type htEntry struct {
	Name       string
	Type       string
	Size       string
	ModifiedAt string
	Link       string
}

func newHTMLTable(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	// compose document.
	doc := &htDoc{}
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
		entry := htEntry{
			Name:       s.GetFirst("name"),
			Type:       s.GetFirst("type"),
			Size:       s.GetFirst("size"),
			ModifiedAt: s.GetFirst("modified_at"),
			Link:       s.GetFirst("link"),
		}
		doc.Entries = append(doc.Entries, entry)
	}
	// execute template.
	buf := new(bytes.Buffer)
	if err := htTmpl1.Execute(buf, doc); err != nil {
		return nil, err
	}
	return ioutil.NopCloser(buf), nil
}

type htDoc2 struct {
	Headers []string
	Index   map[string]int
	Rows    []htRow
}

func (d *htDoc2) initHeader(props []ltsv.Property) {
	d.Headers = make([]string, 0, len(props))
	d.Index = make(map[string]int)
	d.Rows = make([]htRow, 0)
	for _, p := range props {
		d.Index[p.Label] = len(d.Headers)
		d.Headers = append(d.Headers, p.Label)
	}
}

func (d *htDoc2) addRow(props []ltsv.Property) {
	r := htRow {
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

type htRow struct {
	Values []string
	Others string
}

var htTmpl2 = template.Must(template.New("htmltable2").Parse(`<!DOCTYPE html>
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

func newHTMLTable2(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	// compose document.
	doc := &htDoc2{}
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
			doc.initHeader(s.Properties)
			first = false
		}
		doc.addRow(s.Properties)
	}
	// execute template.
	buf := new(bytes.Buffer)
	if err := htTmpl2.Execute(buf, doc); err != nil {
		return nil, err
	}
	return ioutil.NopCloser(buf), nil
}

func init() {
	MustRegister("htmltable", newHTMLTable)
	MustRegister("htmltable2", newHTMLTable2)
}
