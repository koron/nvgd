package ltsv

import (
	"bytes"
	"testing"

	"github.com/koron/nvgd/internal/assert"
)

func TestWriter(t *testing.T) {
	bb := &bytes.Buffer{}
	w := NewWriter(bb, "foo", "bar", "baz")
	w.Write("101", "102", "103")
	w.Write("201", "202", "203")
	w.Write("301", "302", "303")
	assert.Equal(t,
		"foo:101\tbar:102\tbaz:103\n"+
			"foo:201\tbar:202\tbaz:203\n"+
			"foo:301\tbar:302\tbaz:303\n",
		bb.String(), "")
}

func testEscape(t *testing.T, src, want string) {
	t.Helper()
	got := escape(src)
	assert.Equal(t, want, got, "")
}

func TestEscape(t *testing.T) {
	testEscape(t, "noescape", "noescape")
	testEscape(t, `\`, `\\`)
	testEscape(t, "\t", `\t`)
	testEscape(t, "\n", `\n`)
	testEscape(t, "\r", `\r`)
	testEscape(t, "A\tB", `A\tB`)
}

func TestWrite(t *testing.T) {
	bb := &bytes.Buffer{}
	err := Write(bb, []Property{
		{Label: "foo", Value: "123"},
		{Label: "bar", Value: "456"},
		{Label: "baz", Value: "789"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "foo:123\tbar:456\tbaz:789\n", bb.String(), "")
}
