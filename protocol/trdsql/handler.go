// Package trdsql provides TRDSQL's query editor
package trdsql

import (
	"embed"
	"net/url"
	"path"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func init() {
	protocol.MustRegister("trdsql", protocol.ProtocolFunc(Serve))
}

//go:embed assets
var assetFS embed.FS

func Serve(u *url.URL) (*resource.Resource, error) {
	if u.Path == "" {
		u.Path = "/"
		return resource.NewRedirect(u.String()), nil
	}
	reqPath := u.Path
	if reqPath == "/" {
		reqPath = "index.html"
	}
	f, err := assetFS.Open(path.Join("assets", reqPath))
	if err != nil {
		return nil, err
	}
	return resource.New(f).GuessContentType(reqPath), nil
}
