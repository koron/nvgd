package pager

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestPager(t *testing.T) {
	const in = `A
--
B
--
C
--
D
--
E
--
`

	// Single page
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "1"}, in, "A\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "-1"}, in, "E\n--\n")

	// Multiple pages
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "2,3,4"}, in, "B\n--\nC\n--\nD\n--\n")

	// The output order of pages is sorted by page number if they are within
	// the same code.
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "4,3,2"}, in, "B\n--\nC\n--\nD\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "3,4,2"}, in, "B\n--\nC\n--\nD\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "-1,-2"}, in, "D\n--\nE\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "-2,-1"}, in, "D\n--\nE\n--\n")

	// Page range
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "2-4"}, in, "B\n--\nC\n--\nD\n--\n")
	// Reverse range, is sorted.
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "4-2"}, in, "B\n--\nC\n--\nD\n--\n")

	// Duplicated page is ignored
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "1,1"}, in, "A\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "5,-1"}, in, "E\n--\n")

	// If you specify plus and minus for the page number at the same time, the
	// order of the output pages may be changed.
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "pages": "5,-2"}, in, "E\n--\nD\n--\n")

	// "num" option enabled
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "num": "true", "pages": "1"}, in, "(page 1)\nA\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "num": "true", "pages": "-1"}, in, "(page 5)\nE\n--\n")
	filtertest.Check(t, newPager, filter.Params{"eop": "^--", "num": "true", "pages": "1,-1"}, in, "(page 1)\nA\n--\n(page 5)\nE\n--\n")
}

func TestPagerFail(t *testing.T) {
	filtertest.Fail(t, newPager, filter.Params{}, "", `"eop" option is required`)

	// Hard to make regexp.Compile() fail.
	//filtertest.Fail(t, newPager, filter.Params{"eop":"$^"}, "", `invalid "eop" patter: `) 

	filtertest.Fail(t, newPager, filter.Params{"eop": "^--", "pages": "foo"}, "", `invalid "pages": unknown pages item: foo`)
	filtertest.Fail(t, newPager, filter.Params{"eop": "^--", "pages": "0"}, "", `invalid "pages": unknown pages item: 0`)

	// Never happen for current parsePages()
	//filtertest.Fail(t, newPager, filter.Params{"eop": "^--", "pages": ""}, "", `no "pages" choosen`)
}
