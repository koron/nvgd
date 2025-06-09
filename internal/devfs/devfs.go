// Package devfs provides fs.FS of real files for embedded resource development.
package devfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

type devFS struct {
	fs.FS
	prefix  string
	pattern string
}

var _ fs.FS = (*devFS)(nil)

var rootDir string

// Root enables devfs to work when set non-empty path.
func Root(name string) {
	rootDir = name
}

func (fsys *devFS) Open(name string) (fs.File, error) {
	if rootDir == "" {
		return fsys.FS.Open(name)
	}
	if fsys.pattern != "" {
		if ok, _ := filepath.Match(fsys.pattern, name); ok {
			return nil, fs.ErrNotExist
		}
	}
	return os.DirFS(filepath.Join(rootDir, fsys.prefix)).Open(name)
}

// New creates a new fs.FS can be replaced with real file at run time.
func New(base fs.FS, prefix, pattern string) fs.FS {
	return &devFS{
		FS:      base,
		prefix:  prefix,
		pattern: pattern,
	}
}
