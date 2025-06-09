package duckdb

import (
	"strings"
	"testing"

	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/protocol"
)

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "duckdb", protocol.ProtocolFunc(open))
}

func TestIndex(t *testing.T) {
	r := protocoltest.Open(t, "duckdb:///")
	got := protocoltest.ReadAllString(t, r)
	if !strings.Contains(got, "import * as duckdb from ") {
		t.Errorf("not found \"import ... duckdb\" in %+v", got)
	}
}

func TestShowAsView(t *testing.T) {
	r := protocoltest.Open(t, "duckdb:///show-as-view?t=http://127.0.0.1/")
	got := protocoltest.ReadAllString(t, r)
	if !strings.Contains(got, "window.location.replace") {
		t.Errorf("not found \"location.replace\" in %q", got)
	}
}
