package version

import (
	"net/url"

	"github.com/koron/nvgd/internal/version"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func init() {
	protocol.MustRegister("version", &source{})
}

type source struct {
}

func (s *source) Open(u *url.URL) (*resource.Resource, error) {
	return resource.NewString(version.Version), nil
}
