package filter

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/resource"
	"github.com/olekukonko/tablewriter"
)

type doc struct {
	Headers []string
	Index   map[string]int
	Rows    []row

	HasOthers bool
}

func (d *doc) initHeader(props []ltsv.Property) {
	d.Headers = make([]string, 0, len(props)+1)
	d.Index = make(map[string]int)
	d.Rows = make([]row, 0)
	for _, p := range props {
		d.Index[p.Label] = len(d.Headers)
		d.Headers = append(d.Headers, p.Label)
	}
	d.Headers = append(d.Headers, "(others)")
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
	r.Values[len(d.Headers)-1] = r.Others
	d.Rows = append(d.Rows, r)
}

func (d *doc) truncateOthers() {
	d.Headers = d.Headers[:len(d.Headers)-1]
	for i, row := range d.Rows {
		d.Rows[i].Values = row.Values[:len(row.Values)-1]
	}
}

type row struct {
	Values []string
	Others string
}

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
	if !d.HasOthers {
		d.truncateOthers()
	}
	// execute template.
	buf := new(bytes.Buffer)
	t := tablewriter.NewWriter(buf)
	t.SetHeader(d.Headers)
	for _, row := range d.Rows {
		t.Append(row.Values)
	}
	t.Render()
	return r.Wrap(ioutil.NopCloser(buf)), nil
}

func init() {
	filter.MustRegister("texttable", filterFunc)
}
