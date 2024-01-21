package core

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed assets
var assetsFS embed.FS

type resourceServer struct {
	stripFS fs.FS
	fileSrv http.Handler
}

func newResourceServer() (*resourceServer, error) {
	stripFS, err := fs.Sub(assetsFS, "assets")
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
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		name = "."
	}
	f, err := rs.stripFS.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	//fi, err := f.Stat()
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}

func (rs *resourceServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	rs.fileSrv.ServeHTTP(w, r)
}
