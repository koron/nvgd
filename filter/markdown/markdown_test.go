package markdown

import (
	"testing"
)

func TestAppendMarkdownFilter(t *testing.T) {
	for i, tc := range []struct {
		input string
		want  string
	}{
		// No existing query.
		{
			`<a href="doc/foo.md">link</a>`,
			`<a href="doc/foo.md?markdown">link</a>`,
		},
		// No existing query, deeper path.
		{
			`<a href="doc/bar/baz.md">link</a>`,
			`<a href="doc/bar/baz.md?markdown">link</a>`,
		},
		// Existing query string.
		{
			`<a href="doc/foo.md?version=2">link</a>`,
			`<a href="doc/foo.md?version=2&markdown">link</a>`,
		},
		// Existing query with multiple params.
		{
			`<a href="doc/foo.md?lang=en&page=1">link</a>`,
			`<a href="doc/foo.md?lang=en&page=1&markdown">link</a>`,
		},
		// Not a local doc link — should not be modified.
		{
			`<a href="https://example.com/doc/foo.md">link</a>`,
			`<a href="https://example.com/doc/foo.md">link</a>`,
		},
		// No .md extension — should not be modified.
		{
			`<a href="doc/foo.html">link</a>`,
			`<a href="doc/foo.html">link</a>`,
		},
		// Query with special characters in value.
		{
			`<a href="doc/foo.md?name=a%2Fb">link</a>`,
			`<a href="doc/foo.md?name=a%2Fb&markdown">link</a>`,
		},
	} {
		got := string(appendMarkdownFilter([]byte(tc.input)))
		if got != tc.want {
			t.Errorf("#%d: appendMarkdownFilter(%q)\n  got:  %q\n  want: %q", i, tc.input, got, tc.want)
		}
	}
}
