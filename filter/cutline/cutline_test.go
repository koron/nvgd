package cutline

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

const text string = `AAA
BBB
CCC
DDD
EEE
FFF
GGG
HHH
III
JJJ
`

func TestCutline(t *testing.T) {
	filtertest.Check(t, Filter, filter.Params{"start": "DDD", "end": "GGG"}, text, "DDD\nEEE\nFFF\nGGG\n")

	filtertest.Check(t, Filter, filter.Params{"end": "DDD"}, text, "AAA\nBBB\nCCC\nDDD\n")

	filtertest.Check(t, Filter, filter.Params{"start": "GGG"}, text, "GGG\nHHH\nIII\nJJJ\n")

	filtertest.Check(t, Filter, filter.Params{}, text, "AAA\nBBB\nCCC\nDDD\nEEE\nFFF\nGGG\nHHH\nIII\nJJJ\n")
}
