package ltsv

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("not matched:\nactual=%q\nexpected=%q", actual, expected)
	}
}

func TestReader(t *testing.T) {
	r := NewReader(strings.NewReader(
		`foo:123
foo:123	bar:456	baz:789
	    foo:123	bar:456	baz:789		

foo:123	bar:456	foo:789`))
	testRead(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
		},
		Index: map[string][]int{"foo": {0}},
	})
	testRead(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "baz", Value: "789"},
		},
		Index: map[string][]int{"foo": {0}, "bar": {1}, "baz": {2}},
	})
	testRead(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "baz", Value: "789"},
		},
		Index: map[string][]int{"foo": {0}, "bar": {1}, "baz": {2}},
	})
	testRead(t, r, &Set{Index: map[string][]int{}})
	testRead(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "foo", Value: "789"},
		},
		Index: map[string][]int{"foo": {0, 2}, "bar": {1}},
	})
	last, err := r.Read()
	if err != io.EOF {
		t.Fatalf("should be io.EOF for end: %v", err)
	}
	if last != nil {
		t.Errorf("last should be nil: actual=%q", last)
	}
}

func testRead(t *testing.T, r *Reader, expected *Set) {
	actual, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %s\nexpected=%q", err, expected)
	}
	assertEqual(t, actual, expected)
}

func TestGet(t *testing.T) {
	testGet(t, "foo:123", "foo", []string{"123"})
	testGet(t, "foo:123\tfoo:456", "foo", []string{"123", "456"})
	testGet(t, "foo:123\tbar:456\tfoo:789", "foo", []string{"123", "789"})
	testGet(t, "foo:123\tbar:456\tfoo:789", "bar", []string{"456"})
}

func testGet(t *testing.T, src, label string, expected []string) {
	r := NewReader(strings.NewReader(src))
	s, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	actual := s.Get(label)
	assertEqual(t, actual, expected)
}
