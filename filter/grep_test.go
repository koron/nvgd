package filter

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/nvgd/resource"
)

func checkFilter(t *testing.T, f Factory, p Params, in, want string) {
	t.Helper()
	src := resource.NewString(in)
	filter, err := f(src, p)
	if err != nil {
		t.Errorf("failed to create a filter: %s", err)
		return
	}
	b, err := io.ReadAll(filter)
	if err != nil {
		t.Errorf("failed to read from filter: %s", err)
		return
	}
	got := string(b)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected output from filter: -want +got\n%s", diff)
	}
}

func TestGrep(t *testing.T) {
	checkFilter(t, newGrep, Params{"re": "foo"}, "aaa\nfoo\nbbb\n", "foo\n")

	checkFilter(t, newGrep, Params{"re": "foo", "number": "true"},
		"aaa\nfoo\nbbb\n",
		"2: foo\n")
	checkFilter(t, newGrep, Params{"re": "foo", "number": "true"},
		"aaa\nfoo\nbbb\nfoo\nccc\n",
		"2: foo\n4: foo\n")
}

func TestGrepContext(t *testing.T) {
	checkFilter(t, newGrep, Params{"re": "eee", "context": "1"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"ddd\neee\nfff\n")
	checkFilter(t, newGrep, Params{"re": "eee", "context": "2"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"ccc\nddd\neee\nfff\nggg\n")

	checkFilter(t, newGrep, Params{"re": "eee", "context": "1", "number": "true"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"4: ddd\n5: eee\n6: fff\n")

	checkFilter(t, newGrep, Params{"re": "eee", "context": "2", "number": "true"},
		"aaa\nbbb\nccc\nddd\neee\nfff\nggg\nhhh\niii\n",
		"3: ccc\n4: ddd\n5: eee\n6: fff\n7: ggg\n")

	checkFilter(t, newGrep, Params{"re": "XXX", "context": "2"},
		"aaa\nbbb\nccc\nXXX\neee\nXXX\nggg\nhhh\niii\n",
		"bbb\nccc\nXXX\neee\nXXX\nggg\nhhh\n")
}
