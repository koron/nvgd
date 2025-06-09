package core

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/koron/nvgd/internal/devfs"
)

//go:embed assets
var assetsFS embed.FS

type resourceServer struct {
	stripFS fs.FS
	fileSrv http.Handler
}

func newResourceServer() (*resourceServer, error) {
	stripFS, err := fs.Sub(devfs.New(assetsFS, "core", ""), "assets")
	if err != nil {
		return nil, err
	}
	fileSrv := http.FileServer(http.FS(stripFS))
	return &resourceServer{
		stripFS: stripFS,
		fileSrv: fileSrv,
	}, nil
}

func (rs *resourceServer) isServed(name string) (bool, error) {
	name = strings.Trim(name, "/")
	if name == "" {
		name = "."
	}
	f, err := rs.stripFS.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrInvalid) || errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()
	return true, nil
}

func (rs *resourceServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	rs.fileSrv.ServeHTTP(w, r)
}
