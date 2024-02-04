// Package pager provides a NVGD filter to split resource into pages.
package pager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/ringbuf"
	"github.com/koron/nvgd/resource"
)

type Pager struct {
	filter.Base

	rx      *regexp.Regexp
	pages   []int // sorted positive numbers
	lasts   []int // sorted negative numbers
	showNum bool

	pageNum  int
	pageIncr bool

	// fields for lasts enabled
	lastsBuf *ringbuf.Buffer
	lastPage *bytes.Buffer
	lastPut  bool
}

func NewPager(r io.ReadCloser, rx *regexp.Regexp, pages, lasts []int, showNum bool) *Pager {
	f := &Pager{
		rx:       rx,
		pages:    pages,
		lasts:    lasts,
		showNum:  showNum,
		pageNum:  0,
		pageIncr: true,
	}
	if len(lasts) > 0 {
		f.lastsBuf = ringbuf.New(-lasts[0])
		f.lastPage = new(bytes.Buffer)
		f.lastPut = false
	}
	f.Base.Init(r, f.readNext)
	return f
}

func (f *Pager) readNext(buf *bytes.Buffer) error {
	for {
		line, err := f.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			if len(line) == 0 {
				break
			}
		}
		writeNum := false
		if f.pageIncr {
			f.pageNum++
			f.pageIncr = false
			if f.showNum {
				writeNum = true
			}
		}
		// store a line to "lasts"
		if len(f.lasts) > 0 {
			if writeNum {
				fmt.Fprintf(f.lastPage, "(page %d)\n", f.pageNum)
			}
			_, err := f.lastPage.Write(line)
			if err != nil {
				return err
			}
			if !f.lastPut {
				f.lastsBuf.Put(f.lastPage)
				f.lastPut = true
			}
		}
		// check the pattern for end of a page.
		if f.rx.Match(line) {
			f.pageIncr = true
			if len(f.lasts) > 0 {
				f.lastPage = new(bytes.Buffer)
				f.lastPut = false
			}
		}
		// immediate flush
		if contains(f.pages, f.pageNum) {
			if writeNum {
				fmt.Fprintf(buf, "(%d page: nvgd pager)\n", f.pageNum)
			}
			_, err = buf.Write(line)
			return err
		}
	}
	// flush lastsBuf if "lasts" avaiable
	if len(f.lasts) > 0 {
		for i := 0; i < f.lastsBuf.Len(); i++ {
			curr, ok := f.lastsBuf.Peek(i).(*bytes.Buffer)
			if !ok {
				continue
			}
			pnumRel := -f.lastsBuf.Len() + i
			pnumAbs := f.pageNum + 1 + pnumRel
			// don't write if already wrote a page
			if contains(f.lasts, pnumRel) && !contains(f.pages, pnumAbs) {
				if _, err := io.Copy(buf, curr); err != nil {
					return err
				}
			}
		}
	}
	return io.EOF
}

func contains(a []int, v int) bool {
	if len(a) == 0 {
		return false
	}
	mid := len(a) / 2
	if v == a[mid] {
		return true
	}
	if v < a[mid] {
		return contains(a[:mid], v)
	}
	return contains(a[mid+1:], v)
}

var (
	rxCutOne   = regexp.MustCompile(`^-?[1-9]\d*$`)
	rxCutRange = regexp.MustCompile(`^([1-9]\d*)-([1-9]\d*)$`)
)

func parsePages(s string) ([]int, error) {
	pages := make([]int, 0, 2)
	for _, item := range strings.Split(s, ",") {
		if m := rxCutOne.FindString(item); m != "" {
			n, err := strconv.Atoi(m)
			if err != nil {
				return nil, err
			}
			pages = append(pages, n)
		} else if m := rxCutRange.FindStringSubmatch(item); m != nil {
			s, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			e, err := strconv.Atoi(m[2])
			if err != nil {
				return nil, err
			}
			if e > s {
				s, e = e, s
			}
			for i := s; i <= e; i++ {
				pages = append(pages, i)
			}
		} else {
			return nil, fmt.Errorf("unknown pages item: %s", item)
		}
	}
	sort.Ints(pages)
	return pages, nil
}

func findPlus(nums []int) int {
	for i, v := range nums {
		if v > 0 {
			return i
		}
	}
	return len(nums)
}

func newPager(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// parse "eop" as a regexp.
	pat, ok := p["eop"]
	if !ok {
		return nil, errors.New(`"eop" option is required`)
	}
	rx, err := regexp.Compile(pat)
	if err != nil {
		return nil, fmt.Errorf(`invalid "eop" pattern: %w`, err)
	}
	// choose pages
	pages, err := parsePages(p.String("pages", "1"))
	if err != nil {
		return nil, fmt.Errorf(`invalid "pages": %w`, err)
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("no \"pages\" choosen")
	}
	// prepare for last pages
	var lasts []int
	if pages[0] < 0 {
		x := findPlus(pages)
		pages, lasts = pages[x:], pages[:x]
	}
	// "num" shows page number.
	showNum := p.Bool("num", false)
	return r.Wrap(NewPager(r, rx, pages, lasts, showNum)), nil
}

func init() {
	filter.MustRegister("pager", newPager)
}
