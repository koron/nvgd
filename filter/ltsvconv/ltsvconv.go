// Package ltsvconv provides a filter to convert format from LTSV to another.
package ltsvconv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/internal/ltsv"
	"github.com/koron/nvgd/resource"
)

type Format int

const (
	CSV Format = iota + 1
	TSV
)

type wrapWriter struct {
	io.Writer
}

//func (w *wrapWriter) Write(p []byte) (int, error) {
//	return w.w.Write(p)
//}

var _ io.Writer = (*wrapWriter)(nil)

type writer interface {
	Error() error
	Flush()
	Write([]string) error
}

type Converter struct {
	filterbase.Base
	reader *filterbase.LTSVReader

	wrapw *wrapWriter
	roww  writer

	labels []string
}

func New(r io.ReadCloser, format Format) *Converter {
	c := &Converter{
		reader: filterbase.NewLTSVReader(r),
		wrapw:  &wrapWriter{},
	}
	c.Base.Init(r, c.readNext)
	var roww writer
	switch format {
	case CSV:
		roww = csv.NewWriter(c.wrapw)
	case TSV:
		w := csv.NewWriter(c.wrapw)
		w.Comma = '\t'
		roww = w
	}
	c.roww = roww
	return c
}

func (c *Converter) writeRow(w io.Writer, row []string) error {
	c.wrapw.Writer = w
	err := c.roww.Write(row)
	if err != nil {
		return err
	}
	c.roww.Flush()
	return c.roww.Error()
}

func (c *Converter) readNext(buf *bytes.Buffer) error {
	row, err := c.reader.Read()
	if err != nil {
		return err
	}
	if c.labels == nil {
		// Detemine labels with its order, and write the header.
		c.labels = c.setupLabels(row)
		err := c.writeRow(buf, c.labels)
		if err != nil {
			return err
		}
	}
	// Output values according to the order of the labels in the first row.
	values := make([]string, len(c.labels))
	for i, label := range c.labels {
		values[i] = row.GetFirst(label)
	}
	return c.writeRow(buf, values)
}

func (c *Converter) setupLabels(set *ltsv.Set) []string {
	var labels []string
	seen := map[string]struct{}{}
	for _, p := range set.Properties {
		if _, ok := seen[p.Label]; ok {
			continue
		}
		labels = append(labels, p.Label)
	}
	return labels
}

func newLTSVConv(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	var format Format
	formatStr := strings.ToLower(p.String("format", "tsv"))
	switch formatStr {
	case "csv":
		format = CSV
	case "tsv":
		format = TSV
	default:
		return nil, fmt.Errorf("unsupported format: %s", formatStr)
	}
	return r.Wrap(New(r, format)), nil
}

func init() {
	filter.MustRegister("ltsvconv", newLTSVConv)
}
