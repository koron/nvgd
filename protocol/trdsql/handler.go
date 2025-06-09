// Package trdsql provides TRDSQL's query editor
package trdsql

import (
	"embed"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/trdsql", "")

func init() {
	protocol.MustRegister("trdsql", embedresource.New(assetFS))
}
