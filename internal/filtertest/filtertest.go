// Package filtertest provides utilities to test filter.
package filtertest

import (
	"io"
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/assert"
	"github.com/koron/nvgd/resource"
)

// Check checks a filter working or not.
func Check(t *testing.T, f filter.Factory, p filter.Params, in, want string) {
	t.Helper()
	src := resource.NewString(in)
	defer src.Close()
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
	assert.Equal(t, want, string(b), "unexpected output from filter")
}

// Fail checks the filter.Factory will be failed
func Fail(t *testing.T, f filter.Factory, p filter.Params, in, wantErr string) {
	t.Helper()
	src := resource.NewString(in)
	defer src.Close()
	filter, err := f(src, p)
	if err == nil {
		t.Errorf("expected be failed, but succeeded: want=%s", wantErr)
		return
	}
	got := err.Error()
	if got != wantErr {
		t.Errorf("occurred error unmatch:\nwant=%+q\ngot=%+q", wantErr, got)
	}
	if filter != nil {
		t.Errorf("factory failed but returns non-nil filter: %+v", filter)
	}
}
