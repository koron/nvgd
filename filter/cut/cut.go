// Package cut provides "cut" command like filter for NVGD.
package cut

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
)

// Cut represents cut filter.
type Cut struct {
	filter.Base
	delim     []byte
	splitter  SplitFunc
	selectors []cutSelector
	write     cutWriter
}

type SplitFunc func(b []byte) [][]byte

type cutSelector func(dst, src [][]byte) [][]byte
type cutWriter func(io.Writer, []byte) error

// NewCut creates an instance of cut filter.
func NewCut(r io.ReadCloser, delim []byte, selectors []cutSelector, splitFunc SplitFunc) *Cut {
	f := &Cut{
		delim:    delim,
		splitter: splitFunc,
	}
	if len(selectors) == 0 {
		f.write = f.writeAll
	} else {
		f.selectors = selectors
		f.write = f.writeSome
	}
	if f.splitter == nil {
		f.splitter = func(b []byte) [][]byte {
			return bytes.Split(b, delim)
		}
	}
	f.Base.Init(r, f.readNext)
	return f
}

func (f *Cut) readNext(buf *bytes.Buffer) error {
	raw, err := f.ReadLine()
	if err != nil && (err != io.EOF || len(raw) == 0) {
		return err
	}
	return f.write(buf, raw)
}

func (f *Cut) writeAll(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

func (f *Cut) writeSome(w io.Writer, b []byte) error {
	b, lf := splitLF(b)
	src := f.splitter(b)
	var selected [][]byte
	for _, s := range f.selectors {
		selected = s(selected, src)
	}
	_, err := w.Write(bytes.Join(selected, f.delim))
	if err != nil {
		return err
	}
	if len(lf) > 0 {
		_, err := w.Write(lf)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitLF(b []byte) (body, lf []byte) {
	l := len(b)
	if l > 0 && b[l-1] == '\n' {
		l--
		if l > 0 && b[l-1] == '\r' {
			l--
		}
	}
	return b[:l], b[l:]
}

var (
	rxCutOne        = regexp.MustCompile(`^[1-9]\d*$`)
	rxCutRange      = regexp.MustCompile(`^([1-9]\d*)-([1-9]\d*)$`)
	rxCutRangeBegin = regexp.MustCompile(`^([1-9]\d*)-$`)
	rxCutRangeEnd   = regexp.MustCompile(`^-([1-9]\d*)$`)
)

func toCutSelector(s string) ([]cutSelector, error) {
	if s == "" {
		return nil, nil
	}
	var sels []cutSelector
	for _, item := range strings.Split(s, ",") {
		var sel cutSelector
		if m := rxCutOne.FindString(item); m != "" {
			n, err := toCutIndex(m)
			if err != nil {
				return nil, err
			}
			sel = newCutOne(n)
		} else if m := rxCutRange.FindStringSubmatch(item); m != nil {
			s, err := toCutIndex(m[1])
			if err != nil {
				return nil, err
			}
			e, err := toCutIndex(m[2])
			if err != nil {
				return nil, err
			}
			sel = newCutRange(s, e)
		} else if m := rxCutRangeBegin.FindStringSubmatch(item); m != nil {
			n, err := toCutIndex(m[1])
			if err != nil {
				return nil, err
			}
			sel = newCutRangeBegin(n)
		} else if m := rxCutRangeEnd.FindStringSubmatch(item); m != nil {
			n, err := toCutIndex(m[1])
			if err != nil {
				return nil, err
			}
			sel = newCutRangeEnd(n)
		}
		if sel == nil {
			return nil, fmt.Errorf("unknown cut list item: %s", item)
		}
		sels = append(sels, sel)
	}
	return sels, nil
}

func toCutIndex(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n < 1 {
		return 0, fmt.Errorf("small index: %d", n)
	}
	return n - 1, err
}

func newCutOne(value int) cutSelector {
	return func(dst, src [][]byte) [][]byte {
		if value >= len(src) {
			return dst
		}
		return append(dst, src[value])
	}
}

func newCutRange(start, end int) cutSelector {
	if start <= end {
		return func(dst, src [][]byte) [][]byte {
			l := len(src)
			if start >= l {
				return dst
			}
			if end >= l {
				end = l - 1
			}
			return append(dst, src[start:end+1]...)
		}
	}
	return func(dst, src [][]byte) [][]byte {
		l := len(src)
		if end >= l {
			return dst
		}
		if start >= l {
			start = l - 1
		}
		for i := start; i >= end; i-- {
			dst = append(dst, src[i])
		}
		return dst
	}
}

func newCutRangeBegin(n int) cutSelector {
	return func(dst, src [][]byte) [][]byte {
		if n >= len(src) {
			return dst
		}
		return append(dst, src[n:]...)
	}
}

func newCutRangeEnd(n int) cutSelector {
	return func(dst, src [][]byte) [][]byte {
		l := len(src)
		if n >= l {
			n = l - 1
		}
		return append(dst, src[:n+1]...)
	}
}

func splitWhite(s []byte) [][]byte {
	// count separator chunks
	n := 0
	lastWhite := false
	for _, b := range s {
		if b == ' ' || b == '\t' {
			if !lastWhite {
				n++
				lastWhite = true
			}
		} else {
			lastWhite = false
		}
	}
	// split bytes
	a := make([][]byte, 0, n+1)
	lastWhite = false
	x := 0 // start position of a current chunk
	for i, b := range s {
		if b == ' ' || b == '\t' {
			if !lastWhite {
				lastWhite = true
				a = append(a, s[x:i])
			}
		} else {
			if lastWhite {
				lastWhite = false
				x = i
			}
		}
	}
	// complete the last chunk
	if lastWhite {
		a = append(a, s[len(s):])
	} else {
		a = append(a, s[x:])
	}
	return a
}

func newCut(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	delim := []byte(p.String("delim", "\t"))
	var splitter SplitFunc
	if p.Bool("white", false) {
		splitter = splitWhite
	}
	selectors, err := toCutSelector(p.String("list", ""))
	if err != nil {
		return nil, err
	}
	return r.Wrap(NewCut(r, delim, selectors, splitter)), nil
}

func init() {
	filter.MustRegister("cut", newCut)
}
