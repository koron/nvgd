package protocol

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/koron/nvgd/ltsv"
	"github.com/pierrec/lz4"
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
		w   = ltsv.NewWriter(buf, "name", "type", "size", "modified_at", "link")
	)
	path := strings.TrimRight(name, "/")
	// add updir
	if path != "" {
		up := strings.TrimRight(rxLastComponent.ReplaceAllString(path, ""), "/")
		link := fmt.Sprintf("/file://%s/?indexhtml", up)
		err := w.Write("..", "updir", "", "", link)
		if err != nil {
			return nil, err
		}
	}
	for _, fi := range list {
		n := fi.Name()
		var t, link string
		if fi.IsDir() {
			t = "dir"
			link = fmt.Sprintf("/file://%s/%s/?indexhtml", path, n)
		} else {
			t = "file"
			link = fmt.Sprintf("/file://%s/%s", path, n)
		}
		err := w.Write(n, t, strconv.FormatInt(fi.Size(), 10),
			fi.ModTime().Format(time.RFC1123), link)
		if err != nil {
			return nil, err
		}
	}
	return ioutil.NopCloser(buf), nil
}

var (
	rxGz  = regexp.MustCompile(`\.gz$`)
	rxBz2 = regexp.MustCompile(`\.bz2$`)
	rxLz4 = regexp.MustCompile(`\.lz4$`)
)

func (f *File) openFile(name string) (io.ReadCloser, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	// Apply decompress filter.
	if rxGz.MatchString(name) {
		zr, err := gzip.NewReader(r)
		if err != nil {
			r.Close()
			return nil, err
		}
		return zr, nil
	} else if rxBz2.MatchString(name) {
		return newWrapRC(bzip2.NewReader(r), r), nil
	} else if rxLz4.MatchString(name) {
		return newWrapRC(lz4.NewReader(r), r), nil
	}
	return r, nil
}

type wrapRC struct {
	io.Reader
	c io.Closer
}

func newWrapRC(r io.Reader, c io.Closer) io.ReadCloser {
	return &wrapRC{Reader: r, c: c}
}

func (rc *wrapRC) Close() error {
	return rc.c.Close()
}
