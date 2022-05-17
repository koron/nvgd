package file

import (
	"io"
	"testing"
)

func TestBzip2(t *testing.T) {
	r, err := fileOpen("testdata/file_test.bz2")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	s := (string)(b)
	if s != "this is bz2 compressed" {
		t.Errorf("content of \"testdata/file_test.bz2\" is unexpected: %q", s)
	}
}

func TestGzip(t *testing.T) {
	r, err := fileOpen("testdata/file_test.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	s := (string)(b)
	if s != "this is gzip compressed" {
		t.Errorf("content of \"testdata/file_test.gz\" is unexpected: %q", s)
	}
}

func TestLZ4(t *testing.T) {
	r, err := fileOpen("testdata/file_test.lz4")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	s := (string)(b)
	if s != "this is lz4 compressed" {
		t.Errorf("content of \"testdata/file_test.lz4\" is unexpected: %q", s)
	}
}

func TestMultiRC(t *testing.T) {
	f := &File{}
	r, err := f.openMulti([]string{
		"testdata/file_test.bz2",
		"testdata/file_test.gz",
		"testdata/file_test.lz4",
	}, "testdata/;file_test.*")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if s != `this is bz2 compressedthis is gzip compressedthis is lz4 compressed` {
		t.Errorf("multi RC failed: %q", s)
	}
}
