package configp

import (
	"testing"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/assert"
	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/protocol"
)

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "config", protocol.ProtocolFunc(Open))
}

func testContent(t *testing.T, cfg config.Config, want string) {
	Config = cfg
	rsrc := protocoltest.Open(t, "config:///")
	got := protocoltest.ReadAllString(t, rsrc)
	assert.Equal(t, want, got, "")
}

func TestOutput(t *testing.T) {
	testContent(t, config.Config{}, `addr: ""
path_prefix: ""
error_log: ""
access_log: ""
root_contents_file: ""
`)
	testContent(t,
		config.Config{
			Addr:       ":3000",
			PathPrefix: "nvgd",
		},
		`addr: :3000
path_prefix: nvgd
error_log: ""
access_log: ""
root_contents_file: ""
`)
}

func TestHideSecrets(t *testing.T) {
	t.Run("secret_access_key", func(t *testing.T) {
		type secretAccessKeyHolder struct {
			SecretAccessKey config.SecretString `yaml:"secret_access_key"`
		}
		testContent(t,
			config.Config{
				Protocols: config.CustomConfig{
					"secretAccessKeyHolder": secretAccessKeyHolder{
						SecretAccessKey: "should not shown",
					},
				},
			},
			`addr: ""
path_prefix: ""
error_log: ""
access_log: ""
root_contents_file: ""
protocols:
  secretAccessKeyHolder:
    secret_access_key: __SECRET__
`)
	})
}
