// Package examples provides examples resources.
package examples

import (
	"embed"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/examples", "")

func init() {
	protocol.MustRegister("examples", embedresource.New(assetFS))
}
