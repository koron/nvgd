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
