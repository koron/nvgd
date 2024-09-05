package grep

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestGrep(t *testing.T) {
	filtertest.Check(t, newGrep,
		filter.Params{"re": "foo"}, "aaa\nfoo\nbbb\n", "foo\n")

	filtertest.Check(t, newGrep,
		filter.Params{"re": "foo", "number": "true"},
		"aaa\nfoo\nbbb\n",
		"2: foo\n")
	filtertest.Check(t, newGrep,
		filter.Params{"re": "foo", "number": "true"},
		"aaa\nfoo\nbbb\nfoo\nccc\n",
		"2: foo\n4: foo\n")
}

func TestGrepContext(t *testing.T) {
	filtertest.Check(t, newGrep,
		filter.Params{"re": "eee", "context": "1"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"ddd\neee\nfff\n")
	filtertest.Check(t, newGrep,
		filter.Params{"re": "eee", "context": "2"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"ccc\nddd\neee\nfff\nggg\n")

	filtertest.Check(t, newGrep,
		filter.Params{"re": "eee", "context": "1", "number": "true"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"4: ddd\n5: eee\n6: fff\n")

	filtertest.Check(t, newGrep,
		filter.Params{"re": "eee", "context": "2", "number": "true"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"3: ccc\n4: ddd\n5: eee\n6: fff\n7: ggg\n")

	filtertest.Check(t, newGrep,
		filter.Params{"re": "XXX", "context": "2"},
		"aaa\nbbb\nccc\nXXX\neee\nXXX\nggg\nhhh\niii\n",
		"bbb\nccc\nXXX\neee\nXXX\nggg\nhhh\n")
}
