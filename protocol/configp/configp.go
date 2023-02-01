// Package configp provides config protocol for NGVD.
package configp

import (
	"bytes"
	"io"
	"net/url"
	"regexp"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
	"gopkg.in/yaml.v2"
)

// Config hold configuration of nvgd.
var Config config.Config

func init() {
	protocol.MustRegister("config", &ConfigP{})
}

// ConfigP provides config protocol.
type ConfigP struct {
}

var mx = regexp.MustCompile(`(secret_access_key): \S+`)

func (cp *ConfigP) Open(u *url.URL) (*resource.Resource, error) {
	b, err := yaml.Marshal(&Config)
	if err != nil {
		return nil, err
	}
	// FIXME: hide secrets by more generic way.
	b = mx.ReplaceAll(b, []byte("$1: __SECRET__"))
	rs := resource.New(io.NopCloser(bytes.NewReader(b)))
	return rs, nil
}
