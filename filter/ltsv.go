package filter

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/koron/nvgd/resource"
)

type ltsvValue map[string][]string

func parseLTSV(s string) ltsvValue {
	r := ltsvValue{}
	for _, raw := range strings.Split(s, "\t") {
		kv := strings.SplitN(raw, ":", 2)
		if len(kv) != 2 {
			continue
		}
		k, v := kv[0], kv[1]
		slot := r[k]
		r[k] = append(slot, v)
	}
	return r
}

func (v ltsvValue) put(buf *bytes.Buffer) error {
	first := true
	for k, slot := range v {
		for _, v := range slot {
			if first {
				first = false
			} else {
				if _, err := buf.WriteString("\t"); err != nil {
					return err
				}
			}
			if _, err := buf.WriteString(k); err != nil {
				return err
			}
			if _, err := buf.WriteString(":"); err != nil {
				return err
			}
			if _, err := buf.WriteString(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// LTSV represents a structure for LTSV (labeled tab separated value)
type LTSV struct {
	Base
	label string
	re    *regexp.Regexp
	match bool
	cut   []string
}

// NewLTSV creates a new instance of LTSV.
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
	for {
		b, err := l.ReadLine()
		if err != nil {
			return err
		}
		v := parseLTSV(string(b))
		if !l.isMatch(v) {
			continue
		}
		if err := l.filter(v).put(buf); err != nil {
			return err
		}
		if _, err := buf.WriteString("\n"); err != nil {
			return err
		}
	}
}

func (l *LTSV) isMatch(v ltsvValue) bool {
	if l.label == "" {
		return true
	}
	slot, ok := v[l.label]
	if !ok {
		return false
	}
	for _, v := range slot {
		if l.re.MatchString(v) {
			return true
		}
	}
	return false
}

func (l *LTSV) filter(v ltsvValue) ltsvValue {
	if len(l.cut) == 0 {
		return v
	}
	r := ltsvValue{}
	for _, label := range l.cut {
		slot, ok := v[label]
		if !ok {
			continue
		}
		r[label] = slot
	}
	return r
}

func newLTSV(r *resource.Resource, p Params) (*resource.Resource, error) {
	label, re, err := parseGrep(p)
	if err != nil {
		return nil, err
	}
	match := p.Bool("match", true)
	cut := parseCut(p)
	return r.Wrap(NewLTSV(r, label, re, match, cut)), nil
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
