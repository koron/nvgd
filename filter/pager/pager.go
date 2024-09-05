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

	"github.com/koron-go/ringbuf"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

type Pager struct {
	filterbase.Base

	rx      *regexp.Regexp
	pages   []int // sorted positive numbers
	lasts   []int // sorted negative numbers
	showNum bool

	pageNum int

	currPW *pageWriter

	// fields for lasts enabled
	lastsRing *ringbuf.Buffer[*bytes.Buffer]
}

func NewPager(r io.ReadCloser, rx *regexp.Regexp, pages, lasts []int, showNum bool) *Pager {
	f := &Pager{
		rx:      rx,
		pages:   pages,
		lasts:   lasts,
		showNum: showNum,
		pageNum: 0,
	}
	if len(lasts) > 0 {
		f.lastsRing = ringbuf.New[*bytes.Buffer](-lasts[0])
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

		if f.currPW == nil {
			f.pageNum++
			f.currPW = &pageWriter{
				num:     f.pageNum,
				showNum: f.showNum,
				w1:      new(bytes.Buffer),
			}
			if f.lastsRing != nil {
				f.lastsRing.Enqueue(f.currPW.w1)
			}
		}
		hit := contains(f.pages, f.pageNum)
		if hit {
			f.currPW.w2 = buf
		}
		f.currPW.Write(line)
		// Clear pageWriter when "eop" matches. it cause page feeding in next
		// loop.
		if f.rx.Match(line) {
			f.currPW = nil
		}
		// Delegate to another filters when sure to output.
		if hit {
			return nil
		}
	}

	// flush lastsBuf if "lasts" avaiable
	if f.lastsRing != nil {
		for i := 0; i < f.lastsRing.Len(); i++ {
			curr := f.lastsRing.Peek(i)
			pnumRel := -f.lastsRing.Len() + i
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
	rxPageOne   = regexp.MustCompile(`^-?[1-9]\d*$`)
	rxPageRange = regexp.MustCompile(`^([1-9]\d*)-([1-9]\d*)$`)
)

func parsePages(s string) ([]int, error) {
	pages := make([]int, 0, 2)
	for _, item := range strings.Split(s, ",") {
		if m := rxPageOne.FindString(item); m != "" {
			n, err := strconv.Atoi(m)
			if err != nil {
				return nil, err
			}
			pages = append(pages, n)
		} else if m := rxPageRange.FindStringSubmatch(item); m != nil {
			s, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			e, err := strconv.Atoi(m[2])
			if err != nil {
				return nil, err
			}
			if s > e {
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
