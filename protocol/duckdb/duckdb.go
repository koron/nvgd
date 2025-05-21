// Package duckdb provides duckdb protocol for NVGD.
package duckdb

import (
	"embed"

	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var assetFS embed.FS

func init() {
	protocol.MustRegister("duckdb", embedresource.New(assetFS))
}
