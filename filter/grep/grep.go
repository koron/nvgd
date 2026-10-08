// Package grep provides "grep" filter for NVGD.
package grep

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

// Grep represents grep like filter.
type Grep struct {
	filterbase.Base
	reader   *filterbase.LineReader
	currLnum int

	re    *regexp.Regexp
	match bool
	only  bool
	lf    LineFilter
	lnum  bool
	cnum  int

	contextBefore [][]byte
	contextAfter  int
}

// NewGrep creates an instance of grep filter.
func NewGrep(r io.ReadCloser, re *regexp.Regexp, match, only bool, lf LineFilter, lnum bool, cnum int) *Grep {
	g := &Grep{
		reader: filterbase.NewLineReader(r),
		re:     re,
		match:  match,
		only:   only,
		lf:     TrimEOL.Chain(lf),
		lnum:   lnum,
		cnum:   cnum,
	}
	if cnum > 0 {
		g.contextBefore = make([][]byte, 0, cnum)
	}
	g.Base.Init(r, g.readNext)
	return g
}

func (g *Grep) readNext(buf *bytes.Buffer) error {
	for {
		raw, err := g.reader.ReadLine()
		if err != nil && len(raw) == 0 {
			return err
		}
		g.currLnum++
		b := g.lf.Apply(raw)

		if g.only {
			if !g.match {
				continue
			}
			m := g.re.FindAll(b, -1)
			for _, mm := range m {
				err := g.outputln(buf, g.currLnum, mm)
				if err != nil {
					return err
				}
			}
			continue
		}

		if g.re.Match(b) != g.match {
			if g.contextAfter > 0 {
				g.contextAfter--
				return g.output(buf, g.currLnum, raw)
			}
			// add to the before context.
			if g.cnum > 0 {
				if len(g.contextBefore) >= g.cnum {
					copy(g.contextBefore[:g.cnum-1], g.contextBefore[1:g.cnum])
					g.contextBefore = g.contextBefore[:g.cnum-1]
				}
				g.contextBefore = append(g.contextBefore, raw)
			}
			continue
		}
		if g.cnum > 0 {
			// output the before conext.
			for i, d := range g.contextBefore {
				err := g.output(buf, g.currLnum-len(g.contextBefore)+i, d)
				if err != nil {
					return err
				}
			}
			g.contextBefore = g.contextBefore[0:0]
			g.contextAfter = g.cnum
		}
		return g.output(buf, g.currLnum, raw)
	}
}

func (g *Grep) output(buf *bytes.Buffer, lnum int, data []byte) error {
	if g.lnum {
		_, err := fmt.Fprintf(buf, "%d: %s", lnum, data)
		return err
	}
	_, err := buf.Write(data)
	return err
}

func (g *Grep) outputln(buf *bytes.Buffer, lnum int, data []byte) error {
	if g.lnum {
		_, err := fmt.Fprintf(buf, "%d: %s\n", lnum, data)
		return err
	}
	_, err := fmt.Fprintf(buf, "%s\n", data)
	return err
}

func newGrep(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	re, err := regexp.Compile(p.String("re", ""))
	if err != nil {
		return nil, err
	}
	match := p.Bool("match", true)
	only := p.Bool("only", false)
	lnum := p.Bool("number", false)
	cnum := p.Int("context", 0)

	var lf LineFilter
	// field filter
	var (
		field = p.Int("field", 0)
		delim = []byte(p.String("delim", "\t"))
	)
	if field > 0 {
		lf = lf.Chain(NewCutLF(delim, field-1))
	}

	return r.Wrap(NewGrep(r, re, match, only, lf, lnum, cnum)), nil
}

func init() {
	filter.MustRegister("grep", newGrep)
}
