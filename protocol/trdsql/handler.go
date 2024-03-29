// Package trdsql provides TRDSQL's query editor
package trdsql

import (
	"embed"

	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var assetFS embed.FS

func init() {
	protocol.MustRegister("trdsql", embedresource.New(assetFS))
}
