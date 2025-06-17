// Package echarts provides chart editor.
package echarts

import (
	"embed"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "protocol/echarts", "")

func init() {
	p := embedresource.New(assetFS)
	protocol.MustRegister("echarts", p)
}
