// Package assert provides utility to test.
package assert

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal asserts equality of two values.
func Equal(t *testing.T, want, got any, msgfmt string, msgargs ...any) {
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		if msgfmt == "" {
			t.Errorf("not equal: -want +got\n%s", d)
			return
		}
		msg := fmt.Sprintf(msgfmt, msgargs...)
		t.Errorf("not equal: %s: -want +got\n%s", msg, d)
	}
}
