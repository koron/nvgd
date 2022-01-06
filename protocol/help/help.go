// Package help provides help protocol for NVGD.
package help

import (
	"net/url"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func init() {
	protocol.MustRegister("help", &Help{})
}

type Help struct {
}

var Text string

func (hp *Help) Open(u *url.URL) (*resource.Resource, error) {
	return resource.NewString(Text), nil
}
