package ltsv

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestLTSVFilter(t *testing.T) {
	// TODO: fix LTSV filter.
	// internal/ltsvを使うように変更し、その他のテストも書く
	t.Skip("invalid just for now")
	filtertest.Check(t, newLTSV, filter.Params{},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n")
}
