// Package vfs provides virtual file system protocol for NVGD.
// vfs is a protocol that serves compressed files as static content.
package vfs

import (
	"archive/zip"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/httperror"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

type Config struct {
	Archives map[string]string `yaml:"archives,omitempty"`
}

type Fsys struct {
	err error

	vfs vfs.FileSystem
}

func (fsys *Fsys) Open(name string) (*resource.Resource, error) {
	fi, err := fsys.vfs.Stat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return fsys.Open(path.Join(name, "index.html"))
	}
	r, err := fsys.vfs.Open(name)
	if err != nil {
		return nil, err
	}
	return resource.New(r).
			GuessContentType(name).
			Put(resource.SkipFilters, true),
		nil
}

var (
	cfg Config

	mu sync.Mutex

	fsysMap = map[string]*Fsys{}
)

func openFsys(name string) (*Fsys, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	vfs := zipfs.New(r, name)
	return &Fsys{
		vfs: vfs,
	}, nil
}

func getFsys(name string) (*Fsys, error) {
	mu.Lock()
	defer mu.Unlock()
	fsys, ok := fsysMap[name]
	if ok {
		if fsys.err != nil {
			return nil, fsys.err
		}
		return fsys, nil
	}
	fsys, err := openFsys(name)
	if err != nil {
		fsysMap[name] = &Fsys{err: err}
		return nil, err
	}
	fsysMap[name] = fsys
	return fsys, nil

}

func Open(u *url.URL) (*resource.Resource, error) {
	alias := u.Hostname()
	path := u.Path
	name, ok := cfg.Archives[alias]
	if !ok {
		return nil, httperror.Newf(http.StatusNotFound, "no vfs named %q found", alias)
	}
	fsys, err := getFsys(name)
	if err != nil {
		return nil, err
	}
	return fsys.Open(path)
}

func init() {
	protocol.MustRegister("vfs", protocol.ProtocolFunc(Open))
	config.RegisterProtocol("vfs", &cfg)
}
