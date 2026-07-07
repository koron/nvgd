// Package version provides version protocol for NVGD.
package version

import (
	"context"
	"net/url"

	"github.com/koron/nvgd/internal/version"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func init() {
	protocol.MustRegister("version", protocol.ProtocolFunc(Open))
}

func Open(ctx context.Context, u *url.URL) (*resource.Resource, error) {
	return resource.NewString(version.Version), nil
}
