package filter

import (
	"bytes"
	"io"
	"regexp"
)

// Grep represents grep like filter.
type Grep struct {
	Base
	re *regexp.Regexp
}

func NewGrep(r io.ReadCloser, re *regexp.Regexp) *Grep {
	g := &Grep{
		re: re,
	}
	g.Base.Init(r, g.readNext)
	return g
}

func (g *Grep) readNext(buf *bytes.Buffer) error {
	for {
		b, err := g.ReadLine()
		if err != nil {
			return err
		}
		if g.re.Match(b) {
			if _, err := buf.Write(b); err != nil {
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
	return NewGrep(r, re), nil
}

func init() {
	MustRegister("grep", newGrep)
}
