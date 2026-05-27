package vfs

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFsysOpenAndClose(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "test.zip")
	zipContent := "hello from vfs"

	// create a ZIP file
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte(zipContent))
	if err != nil {
		t.Fatal(err)
	}
	err = zw.Close()
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// open via openFsys
	fsys, err := openFsys(zipPath)
	if err != nil {
		t.Fatalf("openFsys failed: %s", err)
	}
	defer fsys.Close()
	if fsys.rc == nil {
		t.Fatal("Fsys.rc is nil, zip ReadCloser not stored")
	}

	// read a file from the VFS — zipfs requires absolute paths
	r, err := fsys.vfs.Open("/hello.txt")
	if err != nil {
		t.Fatalf("vfs.Open failed: %s", err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read failed: %s", err)
	}
	r.Close()
	if string(b) != zipContent {
		t.Fatalf("got %q, want %q", string(b), zipContent)
	}
}

func TestFsysCloseNil(t *testing.T) {
	// Close on a zero-value Fsys must not panic
	fsys := &Fsys{}
	err := fsys.Close()
	if err != nil {
		t.Fatalf("Close on zero Fsys: %s", err)
	}
}
