// Package cutline provides "cutline" filter for NVGD.
//
// This filter is a filter for NVGD that extracts a specific range of lines
// from a text stream. The start and end of the range you want to extract can
// be specified using regular expression patterns, allowing for flexible data
// extraction. For example, this is useful for extracting only specific error
// sections from a log file, or extracting parts of a configuration file.
package cutline

import (
	"bytes"
	"io"
	"regexp"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

func init() {
	filter.MustRegister("cutline", Filter)
}

func Filter(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	var (
		start   = p.String("start", "")
		end     = p.String("end", "")
		rxStart *regexp.Regexp
		rxEnd   *regexp.Regexp
		err     error
	)
	if start != "" {
		if rxStart, err = regexp.Compile(start); err != nil {
			return nil, err
		}
	}
	if end != "" {
		if rxEnd, err = regexp.Compile(end); err != nil {
			return nil, err
		}
	}
	return r.Wrap(New(r, rxStart, rxEnd)), nil
}

type mode int

const (
	before mode = iota
	within
	after
)

type Cutline struct {
	filterbase.Base
	reader *filterbase.LineReader

	start *regexp.Regexp
	end   *regexp.Regexp

	current mode
}

func New(r io.ReadCloser, start, end *regexp.Regexp) *Cutline {
	c := &Cutline{
		reader: filterbase.NewLineReader(r),
		start:  start,
		end:    end,
	}
	if c.start == nil {
		c.current = within
	}
	c.Base.Init(r, c.readNext)
	return c
}

func (c *Cutline) readNext(buf *bytes.Buffer) error {
	for {
		line, err := c.reader.ReadLine()
		if err != nil && len(line) == 0 {
			return err
		}
		switch c.current {
		case before:
			if c.start.Match(line) {
				c.current = within
				return c.output(buf, line)
			}
		case within:
			if c.end != nil && c.end.Match(line) {
				c.current = after
			}
			return c.output(buf, line)
		case after:
			// nothing to do
		}
	}
}

func (c *Cutline) output(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
