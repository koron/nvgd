package help

//go:generate go-assets-builder -p help -o assets.go ../../README.md

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

func (hp *Help) Open(u *url.URL) (*resource.Resource, error) {
	f, err := Assets.Open("/README.md")
	if err != nil {
		return nil, err
	}
	return resource.New(f), nil
}
