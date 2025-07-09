// Package ltsv provides a filter modify (grep and cut) LTSV.
package ltsv

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/internal/ltsv"
	"github.com/koron/nvgd/resource"
)

// LTSV represents a structure for LTSV (labeled tab separated value)
type LTSV struct {
	filterbase.Base
	reader *filterbase.LTSVReader
	label  string
	re     *regexp.Regexp
	match  bool
	cut    []string
}

// NewLTSV creates a new instance of LTSV.
func NewLTSV(r io.ReadCloser, label string, re *regexp.Regexp, match bool, cut []string) *LTSV {
	l := &LTSV{
		label:  label,
		reader: filterbase.NewLTSVReader(r),
		re:     re,
		match:  match,
		cut:    cut,
	}
	l.Base.Init(r, l.readNext)
	return l
}

func (l *LTSV) readNext(buf *bytes.Buffer) error {
	for {
		row, err := l.reader.Read()
		if err != nil {
			return err
		}
		if l.isMatch(row) != l.match {
			continue
		}
		return ltsv.Write(buf, l.filter(row).Properties)
	}
}

func (l *LTSV) isMatch(row *ltsv.Set) bool {
	if l.label == "" {
		return true
	}
	values := row.Get(l.label)
	if len(values) == 0 {
		return false
	}
	for _, v := range values {
		if l.re.MatchString(v) {
			return true
		}
	}
	return false
}

func (l *LTSV) filter(row *ltsv.Set) *ltsv.Set {
	if len(l.cut) == 0 {
		return row
	}
	newRow := ltsv.NewSet()
	for _, label := range l.cut {
		values := row.Get(label)
		if len(values) == 0 {
			continue
		}
		for _, v := range values {
			newRow.Put(label, v)
		}
	}
	return newRow
}

func newLTSV(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	label, re, err := parseGrep(p)
	if err != nil {
		return nil, err
	}
	match := p.Bool("match", true)
	cut := parseCut(p)
	return r.Wrap(NewLTSV(r, label, re, match, cut)), nil
}

func parseGrep(p filter.Params) (label string, pattern *regexp.Regexp, err error) {
	v := strings.SplitN(p.String("grep", ""), ",", 2)
	if len(v) < 2 || v[0] == "" || v[1] == "" {
		return "", nil, nil
	}
	re, err := regexp.Compile(v[1])
	if err != nil {
		return "", nil, err
	}
	return v[0], re, err
}

func parseCut(p filter.Params) []string {
	s := p.String("cut", "")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func init() {
	filter.MustRegister("lstv", newLTSV)
}
