// Package duckdb provides duckdb protocol for NVGD.
package duckdb

import (
	"embed"
	"log"

	"github.com/koron/nvgd/internal/templateresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var assetFS embed.FS

const Version = "1.29.1-dev132.0"

func init() {
	r, err := templateresource.New(assetFS,
		templateresource.WithConstant(map[string]any{
			"version": Version,
		}))
	if err != nil {
		log.Fatalf("failed to initialize protocol/duckdb: %s", err)
	}
	protocol.MustRegister("duckdb", r)
}
