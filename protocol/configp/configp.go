// Package configp provides config protocol for NGVD.
package configp

import (
	"bytes"
	"io"
	"net/url"

	"github.com/goccy/go-yaml"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

// Config hold configuration of nvgd.
var Config config.Config

func init() {
	protocol.MustRegister("config", protocol.ProtocolFunc(Open))
}

func Open(u *url.URL) (*resource.Resource, error) {
	b, err := yaml.Marshal(&Config)
	if err != nil {
		return nil, err
	}
	rs := resource.New(io.NopCloser(bytes.NewReader(b)))
	return rs, nil
}
