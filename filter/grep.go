package filter

import (
	"bytes"
	"io"
	"regexp"
)

// Grep represents grep like filter.
type Grep struct {
	Base
	re    *regexp.Regexp
	match bool
	lf    LineFilter
}

// NewGrep creates an instance of grep filter.
func NewGrep(r io.ReadCloser, re *regexp.Regexp, match bool, lf LineFilter) *Grep {
	g := &Grep{
		re:    re,
		match: match,
		lf:    TrimEOL.Chain(lf),
	}
	g.Base.Init(r, g.readNext)
	return g
}

func (g *Grep) readNext(buf *bytes.Buffer) error {
	for {
		raw, err := g.ReadLine()
		if err != nil {
			return err
		}
		b := g.lf.Apply(raw)
		if g.re.Match(b) == g.match {
			if _, err := buf.Write(raw); err != nil {
				return err
			}
			return nil
		}
	}
}

func newGrep(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	re, err := regexp.Compile(p.String("re", ""))
	if err != nil {
		return nil, err
	}
	match := p.Bool("match", true)
	var lf LineFilter
	// field filter
	var (
		field = p.Int("field", 0)
		delim = []byte(p.String("delim", "\t"))
	)
	if field > 0 {
		lf = lf.Chain(NewCutLF(delim, field-1))
	}
	return NewGrep(r, re, match, lf), nil
}

func init() {
	MustRegister("grep", newGrep)
}
