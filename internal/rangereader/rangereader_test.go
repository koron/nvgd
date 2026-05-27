package rangereader

import (
	"io"
	"strings"
	"testing"
)

func TestNonCloser(t *testing.T) {
	// strings.Reader does not implement io.Closer —
	// this must not panic.
	_, err := New(strings.NewReader("hello world"), 0, 10, 100)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestReadCloserImplements(t *testing.T) {
	var _ io.ReadCloser = (*RangeReader)(nil)
}
