package jsonarray

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestJSONArray(t *testing.T) {
	filtertest.Check(t, newFilter, filter.Params{}, "aaa\nbbb\nccc", `["aaa",
"bbb",
"ccc"]
`)
	filtertest.Check(t, newFilter, filter.Params{}, "aaa\nbbb\nccc\n", `["aaa",
"bbb",
"ccc",
""]
`)
}
