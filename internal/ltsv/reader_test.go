package ltsv

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/koron/nvgd/internal/assert"
)

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
	assert.Equal(t, expected, actual, "")
}

func TestLongLine(t *testing.T) {
	// This file (long_line.ltsv) is a one-line LTSV file with 10 properties,
	// each with a label and value of a number from 0 to 9.
	//
	// Each property label is a one-character string corresponding to the
	// corresponding number. Each property value is a 500-character string
	// corresponding to the corresponding number. It is important that the
	// result exceeds 4096 bytes in one line.
	f, err := os.Open(filepath.Join("testdata", "longline_001.ltsv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReaderSize(f, 1*1024*1024)
	set, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	// Verify the set contents.
	if got := len(set.Properties); got != 10 {
		t.Errorf("unexpected size to read: want=%d got=%d", 10, got)
	}
	for i, p := range set.Properties {
		wantLabel := strconv.Itoa(i)
		gotLabel := p.Label
		assert.Equal(t, wantLabel, gotLabel, "incorrect label for #%d entry", i)
		wantValue := strings.Repeat(wantLabel, 500)
		assert.Equal(t, wantValue, p.Value, "incorrect value for #%d entry", i)
	}
}

func TestLongLine2(t *testing.T) {
	// Test to confirm that the line can be read twice as long. Confirming that
	// the buffer extension occurs two or more times.
	f, err := os.Open(filepath.Join("testdata", "longline_002.ltsv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReaderSize(f, 1*1024*1024)
	set, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	// Verify the set contents.
	if got := len(set.Properties); got != 20 {
		t.Errorf("unexpected size to read: want=%d got=%d", 20, got)
	}
	for i, p := range set.Properties {
		wantLabel := strconv.Itoa(i % 10)
		gotLabel := p.Label
		assert.Equal(t, wantLabel, gotLabel, "incorrect label for #%d entry", i)
		wantValue := strings.Repeat(wantLabel, 500)
		assert.Equal(t, wantValue, p.Value, "incorrect value for #%d entry", i)
	}
}
