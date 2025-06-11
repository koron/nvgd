// Package opfs provides a UI operates with OPFS.
package opfs

import (
	"embed"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/opfs", "")

func init() {
	p := embedresource.New(assetFS, embedresource.WithSkipFilter())
	protocol.MustRegister("opfs", p)
}
