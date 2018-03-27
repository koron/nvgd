package configp

import (
	"bytes"
	"io/ioutil"
	"net/url"

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

func (cp *ConfigP) Open(u *url.URL) (*resource.Resource, error) {
	b, err := yaml.Marshal(&Config)
	if err != nil {
		return nil, err
	}
	rs := resource.New(ioutil.NopCloser(bytes.NewReader(b)))
	return rs, nil
}
