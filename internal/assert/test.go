// Package assert provides utility to test.
package assert

import (
	"fmt"
	"reflect"
	"testing"
)

// Equals asserts equality of two values.
func Equals(t *testing.T, actual, expected interface{}, format string, a ...interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nactual=%+v\nexpected=%+v", msg, actual, expected)
	}
}
