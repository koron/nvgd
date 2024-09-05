package count

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestCount(t *testing.T) {
	filtertest.Check(t, newCount, filter.Params{}, "a\nb\nc\nd\n", "4")
	filtertest.Check(t, newCount, filter.Params{}, "a\nb\nc\nd", "4")
	filtertest.Check(t, newCount, filter.Params{}, "a\nb\nc\n", "3")
}
