package filter

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

type LTSV struct {
	Base
	label string
	re    *regexp.Regexp
	match bool
	cut   []string
}

func NewLTSV(r io.ReadCloser, label string, re *regexp.Regexp, match bool, cut []string) *LTSV {
	l := &LTSV{
		label: label,
		re:    re,
		match: match,
		cut:   cut,
	}
	l.Base.Init(r, l.readNext)
	return l
}

func (l *LTSV) readNext(buf *bytes.Buffer) error {
	// TODO:
	return io.EOF
}

func newLTSV(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	label, re, err := parseGrep(p)
	if err != nil {
		return nil, err
	}
	match := p.Bool("match", true)
	cut := parseCut(p)
	return NewLTSV(r, label, re, match, cut), nil
}

func parseGrep(p Params) (label string, pattern *regexp.Regexp, err error) {
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

func parseCut(p Params) []string {
	s := p.String("cut", "")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func init() {
	MustRegister("lstv", newLTSV)
}
