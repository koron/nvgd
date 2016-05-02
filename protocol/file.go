package protocol

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// File is file protocol handler.
type File struct {
}

func init() {
	MustRegister("file", &File{})
}

// Open opens a URL as file.
func (f *File) Open(u *url.URL) (io.ReadCloser, error) {
	// TODO: consider relative path.
	name := u.Path
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.openDir(name)
	}
	return f.openFile(name)
}

func (f *File) openDir(name string) (io.ReadCloser, error) {
	list, err := ioutil.ReadDir(name)
	if err != nil {
		log.Printf("ReadDir failed: %s", name)
		return nil, err
	}
	var (
		buf = &bytes.Buffer{}
		out = make([]string, 0, 4)
	)
	for _, fi := range list {
		t := "file"
		if fi.IsDir() {
			t = "dir"
		}
		out := append(out, fi.Name(), t, strconv.FormatInt(fi.Size(), 10), fi.ModTime().Format(time.RFC1123))
		_, err := buf.WriteString(strings.Join(out, "\t") + "\n")
		if err != nil {
			return nil, err
		}
		out = out[0:0]
	}
	return ioutil.NopCloser(buf), nil
}

func (f *File) openFile(name string) (io.ReadCloser, error) {
	return os.Open(name)
}
