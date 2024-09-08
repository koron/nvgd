package version

import (
	"testing"

	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/internal/version"
	"github.com/koron/nvgd/protocol"
)

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "version", protocol.ProtocolFunc(Open))
}

func TestVersion(t *testing.T) {
	rsrc := protocoltest.Open(t, "version:")
	got := protocoltest.ReadAllString(t, rsrc)
	want := version.Version
	if got != want {
		t.Errorf("wrong versoin: want=%q got=%q", want, got)
	}
}
