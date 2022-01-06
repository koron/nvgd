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
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/koron/nvgd/common_const"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/ltsv"
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
	name := u.Path
	m, err := filepath.Glob(name)
	if err != nil {
		return nil, err
	}
	if len(m) == 1 {
		return f.openOne(name)
	}
	return f.openMulti(m, name)
}

func (f *File) openOne(name string) (*resource.Resource, error) {
	// TODO: consider relative path.
	if !fc.isAccessible(name) {
		return nil, fmt.Errorf("forbidden: %s", name)
	}
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return fileOpenDir(name)
	}
	rc, err := fileOpen(name)
	if err != nil {
		return nil, err
	}
	return resource.New(rc), nil
}

func (f *File) openMulti(names []string, pattern string) (*resource.Resource, error) {
	readers := make([]io.Reader, 0, len(names))
	for _, n := range names {
		if !fc.isAccessible(n) {
			continue
		}
		readers = append(readers, newDelayFile(n))
	}
	if len(readers) == 0 {
		return nil, fmt.Errorf("no matches: %s", pattern)
	}
	return resource.New(newMultiRC(readers...)), nil
}

func fileOpenDir(name string) (*resource.Resource, error) {
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
	for _, fi0 := range list {
		n := fi0.Name()
		fi, err := os.Stat(filepath.Join(path, n))
		if err != nil {
			return nil, err
		}
		var t, link, download string
		if fi.IsDir() {
			t = "dir"
			link = fmt.Sprintf("/file://%s/%s/?indexhtml", path, n)
		} else {
			t = "file"
			link = fmt.Sprintf("/file://%s/%s", path, n)
			download = link + "?all&download"
		}
		err = w.Write(n, t, strconv.FormatInt(fi.Size(), 10),
			fi.ModTime().Format(time.RFC1123), link, download)
		if err != nil {
			return nil, err
		}
	}
	rs := resource.New(ioutil.NopCloser(buf))
	rs.Put(common_const.LTSV, true)
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

func fileOpen(name string) (io.ReadCloser, error) {
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

func fileFailure(err error, readers []io.Reader) error {
	for _, r := range readers {
		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}
	}
	return err
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

type delayFile struct {
	rc  io.ReadCloser
	n   string
	err error
}

func newDelayFile(name string) *delayFile {
	return &delayFile{n: name}
}

func (d *delayFile) Read(b []byte) (int, error) {
	if d.err != nil {
		return 0, d.err
	}
	if d.rc == nil {
		d.rc, d.err = fileOpen(d.n)
		if d.err != nil {
			return 0, d.err
		}
	}
	return d.rc.Read(b)
}

func (d *delayFile) Close() error {
	if d.err != nil {
		return d.err
	}
	if d.rc == nil {
		d.err = io.EOF
		return nil
	}
	d.err = d.rc.Close()
	d.rc = nil
	return d.err
}

type multiRC struct {
	io.Reader
	rcs []io.ReadCloser
}

func newMultiRC(readers ...io.Reader) *multiRC {
	rcs := make([]io.ReadCloser, 0, len(readers))
	for _, r := range readers {
		if rc, ok := r.(io.ReadCloser); ok {
			rcs = append(rcs, rc)
		}
	}
	return &multiRC{
		Reader: io.MultiReader(readers...),
		rcs:    rcs,
	}
}

func (mrc *multiRC) Close() error {
	for _, rc := range mrc.rcs {
		rc.Close()
	}
	mrc.rcs = nil
	return nil
}
