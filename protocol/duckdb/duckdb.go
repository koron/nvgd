// Package duckdb provides duckdb protocol for NVGD.
package duckdb

import (
	"embed"
	"io/fs"
	"net/url"
	"sync"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/templateresource"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/duckdb", "")

// Version of duckdb/duckdb-wasm.
// See https://cdn.jsdelivr.net/npm/@duckdb/duckdb-wasm/ for newer version.
const Version = "1.29.1-dev263.0"

var getRsrc = sync.OnceValues(func() (*templateresource.TemplateResource, error) {
	fsys, err := fs.Sub(assetFS, "assets")
	if err != nil {
		return nil, err
	}
	return templateresource.New(fsys, templateresource.WithConstant(map[string]any{
		"version": Version,
	}))
})

func init() {
	protocol.MustRegister("duckdb", protocol.ProtocolFunc(open))
}

func open(u *url.URL) (*resource.Resource, error) {
	tmplRsrc, err := getRsrc()
	if err != nil {
		return nil, err
	}
	r, err := tmplRsrc.Open(u)
	if err != nil {
		return nil, err
	}
	r.Put(resource.SkipFilters, true)
	return r, nil
}
