package ltsv

import (
	"reflect"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	r := NewReader(strings.NewReader(
		`foo:123
foo:123	bar:456	baz:789
	    foo:123	bar:456	baz:789		
foo:123	bar:456	foo:789
`))
	testASet(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
		},
		Index: map[string][]int{"foo": {0}},
	})
	testASet(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "baz", Value: "789"},
		},
		Index: map[string][]int{"foo": {0}, "bar": {1}, "baz": {2}},
	})
	testASet(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "baz", Value: "789"},
		},
		Index: map[string][]int{"foo": {0}, "bar": {1}, "baz": {2}},
	})
	testASet(t, r, &Set{
		Properties: []Property{
			{Label: "foo", Value: "123"},
			{Label: "bar", Value: "456"},
			{Label: "foo", Value: "789"},
		},
		Index: map[string][]int{"foo": {0, 2}, "bar": {1}},
	})
}

func testASet(t *testing.T, r *Reader, expected *Set) {
	actual, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("not matched:\nactual=%q\nexpected=%q", actual, expected)
	}
}
