package ltsvconv

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestLTSVConv(t *testing.T) {
	filtertest.Check(t, newLTSVConv,
		filter.Params{},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a\tb\tc\n111\t222\t333\n444\t555\t666\n777\t888\t999\n")

	filtertest.Check(t, newLTSVConv,
		filter.Params{"format": "tsv"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a\tb\tc\n111\t222\t333\n444\t555\t666\n777\t888\t999\n")

	filtertest.Check(t, newLTSVConv,
		filter.Params{"format": "csv"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a,b,c\n111,222,333\n444,555,666\n777,888,999\n")
}
