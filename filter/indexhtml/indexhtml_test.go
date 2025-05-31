package indexhtml

import (
	"testing"

	"github.com/koron/nvgd/config"
)

func TestPathPrefix(t *testing.T) {
	for i, c := range []struct{ prefix, path, want string }{
		{"", "/foo/bar", "/foo/bar"},

		{"/qux/", "/foo/bar", "/qux/foo/bar"},
		{"/qux", "/foo/bar", "/qux/foo/bar"},
		{"/qux", "foo/bar", "/qux/foo/bar"},
		{"/qux/", "foo/bar", "/qux/foo/bar"},

		{"/qux/", "/foo://bar", "/qux/foo://bar"},
	} {
		config.Root().PathPrefix = c.prefix
		got := pathPrefix(c.path)
		if got != c.want {
			t.Errorf("unexpected at %d %+v:\ngot=%s", i, c, got)
		}
	}
}
