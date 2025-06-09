// Package duckdb provides duckdb protocol for NVGD.
package duckdb

import (
	"embed"
	"log"
	"net/url"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/templateresource"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/duckdb", "")

const Version = "1.29.1-dev208.0"

var tmplRsrc *templateresource.TemplateResource

func init() {
	// TODO: Package templateresource loads the file at the time of this
	// init(), and even if you switch it with devfs as a startup option, it
	// will not be reflected. I would like a mechanism to delay loading and
	// discard the cache when updating.
	r, err := templateresource.New(assetFS,
		templateresource.WithConstant(map[string]any{
			"version": Version,
		}))
	if err != nil {
		log.Fatalf("failed to initialize protocol/duckdb: %s", err)
	}
	tmplRsrc = r
	protocol.MustRegister("duckdb", protocol.ProtocolFunc(open))
}

func open(u *url.URL) (*resource.Resource, error) {
	r, err := tmplRsrc.Open(u)
	if err != nil {
		return nil, err
	}
	r.Put(resource.SkipFilters, true)
	return r, nil
}
