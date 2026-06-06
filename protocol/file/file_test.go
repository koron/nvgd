package file

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/koron/nvgd/internal/assert"
)

func TestBzip2(t *testing.T) {
	r, stripped, err := fileOpen("testdata/file_test.bz2", false)
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
	assert.Equal(t, "testdata/file_test", stripped, "unmatch stripped")
}

func TestBzip2KeepCompress(t *testing.T) {
	r, stripped, err := fileOpen("testdata/file_test.bz2", true)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/file_test.bz2")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got, "")

	assert.Equal(t, "testdata/file_test.bz2", stripped, "should not stripped")
}

func TestGzip(t *testing.T) {
	r, stripped, err := fileOpen("testdata/file_test.gz", false)
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
	assert.Equal(t, "testdata/file_test", stripped, "unmatch stripped")
}

func TestLZ4(t *testing.T) {
	r, stripped, err := fileOpen("testdata/file_test.lz4", false)
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
	assert.Equal(t, "testdata/file_test", stripped, "unmatch stripped")
}

func TestActualOpenGlobNoMatch(t *testing.T) {
	f := &File{}
	_, err := f.actualOpen(&url.URL{Path: "testdata/nonexistent*"}, false)
	if err == nil {
		t.Fatal("expected error for glob with no matches")
	}
	if !strings.HasPrefix(err.Error(), "no matches:") {
		t.Errorf("expected 'no matches:' error, got: %v", err)
	}
}

type errReader struct {
	io.Reader
}

func (r *errReader) Close() error {
	return fmt.Errorf("close error")
}

func TestMultiRCCollectsErrors(t *testing.T) {
	r1 := struct {
		io.Reader
		io.Closer
	}{bytes.NewReader([]byte("a")), io.NopCloser(nil)}
	r2 := &errReader{Reader: bytes.NewReader([]byte("b"))}
	mrc := newMultiRC(r1, r2)
	err := mrc.Close()
	if err == nil {
		t.Fatal("expected error from multiRC.Close()")
	}
	if !strings.Contains(err.Error(), "close error") {
		t.Errorf("expected 'close error' in: %v", err)
	}
}

func TestMultiRC(t *testing.T) {
	f := &File{}
	r, err := f.openMulti([]string{
		"testdata/file_test.bz2",
		"testdata/file_test.gz",
		"testdata/file_test.lz4",
	}, "testdata/;file_test.*", false)
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
