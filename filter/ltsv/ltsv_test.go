package ltsv

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestLTSVFilter(t *testing.T) {
	filtertest.Check(t, newLTSV,
		filter.Params{},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n")

	filtertest.Check(t, newLTSV,
		filter.Params{"grep": "b,555"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:444\tb:555\tc:666\n")
	filtertest.Check(t, newLTSV,
		filter.Params{"grep": "c,333"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:111\tb:222\tc:333\n")
	filtertest.Check(t, newLTSV,
		filter.Params{"grep": "a,777"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:777\tb:888\tc:999\n")

	filtertest.Check(t, newLTSV,
		filter.Params{"match": "false", "grep": "b,555"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:111\tb:222\tc:333\na:777\tb:888\tc:999\n")

	filtertest.Check(t, newLTSV,
		filter.Params{"cut": "a"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"a:111\na:444\na:777\n")
	filtertest.Check(t, newLTSV,
		filter.Params{"cut": "c,b"},
		"a:111\tb:222\tc:333\na:444\tb:555\tc:666\na:777\tb:888\tc:999\n",
		"c:333\tb:222\nc:666\tb:555\nc:999\tb:888\n")
}
