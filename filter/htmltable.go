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
<table>
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

func init() {
	MustRegister("htmltable", newHTMLTable)
}
