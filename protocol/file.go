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

	"github.com/koron/nvgd/common_const"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/resource"
	"github.com/pierrec/lz4"
)

// File is file protocol handler.
type File struct {
}

// FileConfig provides configuration for file protocol.
type FileConfig struct {
	// Locations is allowed paths to access.
	Locations []string `yaml:"locations"`

	// Forbiddens is fobidden paths to access. It overrides Locations.
	Forbiddens []string `yaml:"forbiddens"`
}

func match(path string, paths []string, defaultValue bool) bool {
	if len(paths) == 0 {
		return defaultValue
	}
	for _, s := range paths {
		if strings.HasPrefix(path, s) {
			return true
		}
	}
	return false
}

func (fc FileConfig) isAccessible(path string) bool {
	return match(path, fc.Locations, true) &&
		!match(path, fc.Forbiddens, false)
}

var fc FileConfig

func init() {
	MustRegister("file", &File{})
	config.RegisterProtocol("file", &fc)
}

// Open opens a URL as file.
func (f *File) Open(u *url.URL) (*resource.Resource, error) {
	// TODO: consider relative path.
	name := u.Path
	if !fc.isAccessible(name) {
		return nil, fmt.Errorf("forbidden: %s", name)
	}
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.openDir(name)
	}
	rc, err := f.openFile(name)
	if err != nil {
		return nil, err
	}
	return resource.New(rc), nil
}

func (f *File) openDir(name string) (*resource.Resource, error) {
	list, err := ioutil.ReadDir(name)
	if err != nil {
		log.Printf("ReadDir failed: %s", name)
		return nil, err
	}
	var (
		buf = &bytes.Buffer{}
		w   = ltsv.NewWriter(buf, "name", "type", "size", "modified_at", "link", "download")
	)
	path := strings.TrimRight(name, "/")
	for _, fi := range list {
		n := fi.Name()
		var t, link, download string
		if fi.IsDir() {
			t = "dir"
			link = fmt.Sprintf("/file://%s/%s/?indexhtml", path, n)
		} else {
			t = "file"
			link = fmt.Sprintf("/file://%s/%s", path, n)
			download = link + "?download"
		}
		err := w.Write(n, t, strconv.FormatInt(fi.Size(), 10),
			fi.ModTime().Format(time.RFC1123), link, download)
		if err != nil {
			return nil, err
		}
	}
	rs := resource.New(ioutil.NopCloser(buf))
	rs.Put(Small, true)
	// add updir
	if path != "" {
		up := strings.TrimRight(rxLastComponent.ReplaceAllString(path, ""), "/")
		link := fmt.Sprintf("/file://%s/?indexhtml", up)
		rs.Put(common_const.UpLink, link)
	}
	return rs, nil
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
