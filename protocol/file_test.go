package protocol

import (
	"io/ioutil"
	"testing"
)

func TestBzip2(t *testing.T) {
	r, err := fileOpen("testdata/file_test.bz2")
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r)
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
	b, err := ioutil.ReadAll(r)
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
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	s := (string)(b)
	if s != "this is lz4 compressed" {
		t.Errorf("content of \"testdata/file_test.lz4\" is unexpected: %q", s)
	}
}
