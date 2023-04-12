// Package examples provides examples resources.
package examples

import (
	"embed"

	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var assetFS embed.FS

func init() {
	protocol.MustRegister("examples", embedresource.New(assetFS))
}
