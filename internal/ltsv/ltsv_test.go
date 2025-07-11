package ltsv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertEqual(t *testing.T, want, got any) {
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("not matched: -want +got\n%s", d)
	}
}

func TestGet(t *testing.T) {
	testGet(t, []Property{{"foo", "123"}}, "foo", []string{"123"})
	testGet(t, []Property{{"foo", "123"}, {"foo", "456"}},
		"foo", []string{"123", "456"})

	data := []Property{
		{"foo", "123"},
		{"bar", "456"},
		{"foo", "789"},
	}
	testGet(t, data, "foo", []string{"123", "789"})
	testGet(t, data, "bar", []string{"456"})

	testGet(t, []Property{{"foo", "123"}}, "bar", nil)
}

func testGet(t *testing.T, props []Property, label string, want []string) {
	t.Helper()
	s := NewSet()
	s.PutProperties(props)
	got := s.Get(label)
	assertEqual(t, got, want)
}

func TestNewSet(t *testing.T) {
	p := NewSet()
	if p == nil {
		t.Fatal("NewSet returns nil")
	}

	if !p.Empty() {
		t.Error("Empty should true for a new Set")
	}
	p.Put("foo", "bar")
	if p.Empty() {
		t.Error("Empty should false after added some props")
	}
}

func TestGetFirst(t *testing.T) {
	p := NewSet()
	p.Put("foo", "first")
	p.Put("foo", "second")
	if got := len(p.Properties); got != 2 {
		t.Errorf("count of properties should be 2: got=%d", got)
	}
	assertEqual(t, 2, len(p.Properties))
	assertEqual(t, 1, len(p.Index))
	assertEqual(t, "first", p.GetFirst("foo"))

	// GetFirst with unexist key.
	assertEqual(t, "", p.GetFirst("bar"))
}
