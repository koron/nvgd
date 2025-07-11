package cut

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/assert"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestCutSelector(t *testing.T) {
	const src = "A\tB\tC\tD\tE\tF\tG\tH\tI\tJ\tK\tL\tM\tN\tO\tP\tQ\tR\tS\tT\tU\tV\tW\tX\tY\tZ\na\tb\tc\td\te\tf\tg\th\ti\tj\tk\tl\tm\tn\to\tp\tq\tr\ts\tt\tu\tv\tw\tx\ty\tz"

	// single selection
	filtertest.Check(t, newCut, filter.Params{"list": "1"}, src, "A\na")
	filtertest.Check(t, newCut, filter.Params{"list": "12"}, src, "L\nl")

	// range
	filtertest.Check(t, newCut, filter.Params{"list": "11-15"}, src, "K\tL\tM\tN\tO\nk\tl\tm\tn\to")
	filtertest.Check(t, newCut, filter.Params{"list": "21-25"}, src, "U\tV\tW\tX\tY\nu\tv\tw\tx\ty")

	// reverse range
	filtertest.Check(t, newCut, filter.Params{"list": "13-11"}, src, "M\tL\tK\nm\tl\tk")

	// open range
	filtertest.Check(t, newCut, filter.Params{"list": "24-"}, src, "X\tY\tZ\nx\ty\tz")
	filtertest.Check(t, newCut, filter.Params{"list": "-3"}, src, "A\tB\tC\na\tb\tc")

	// combinations
	filtertest.Check(t, newCut, filter.Params{"list": "1,5"}, src, "A\tE\na\te")
	filtertest.Check(t, newCut, filter.Params{"list": "12,26"}, src, "L\tZ\nl\tz")

	// TODO: need complex combinations
}

func TestCutEmpty(t *testing.T) {
	// empty lines are at top, middle and bottom.
	const src = `
A1	B1	C1
A2	B2	C2

A4	B4	C4
A5	B5	C5

`
	filtertest.Check(t, newCut, filter.Params{"list": "1"}, src, "\nA1\nA2\n\nA4\nA5\n\n")
	filtertest.Check(t, newCut, filter.Params{"list": "2"}, src, "\nB1\nB2\n\nB4\nB5\n\n")
	filtertest.Check(t, newCut, filter.Params{"list": "3"}, src, "\nC1\nC2\n\nC4\nC5\n\n")
}

func TestSplitWhite(t *testing.T) {
	for i, c := range []struct {
		data string
		want []string
	}{
		{"", []string{""}},

		{"foo", []string{"foo"}},
		{" foo ", []string{"", "foo", ""}},
		{"foo b", []string{"foo", "b"}},

		{"foo bar\tbaz", []string{"foo", "bar", "baz"}},
		{"foo bar\tbaz", []string{"foo", "bar", "baz"}},
		{"foo  \t\t  bar\t\t  \t\tbaz", []string{"foo", "bar", "baz"}},
	} {
		raw := splitWhite([]byte(c.data))
		got := make([]string, 0, len(raw))
		for _, v := range raw {
			got = append(got, string(v))
		}
		assert.Equal(t, c.want, got, "case #%d %+v failed", i, c)
	}
}

func TestCutWhite(t *testing.T) {
	const src = `Jan 31
Feb  1
Feb 23
`
	filtertest.Check(t, newCut, filter.Params{"white": "true", "list": "1"},
		src, "Jan\nFeb\nFeb\n")
	filtertest.Check(t, newCut, filter.Params{"white": "true", "list": "2"},
		src, "31\n1\n23\n")
}
