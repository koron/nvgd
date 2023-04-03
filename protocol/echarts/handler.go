// Package echarts provides chart editor.
package echarts

import (
	"embed"

	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/protocol"
)

//go:embed assets
var assetFS embed.FS

func init() {
	protocol.MustRegister("echarts", embedresource.New(assetFS))
}
