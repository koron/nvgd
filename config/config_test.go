package config

import (
	"path/filepath"
	"testing"

	"github.com/koron/nvgd/internal/assert"
)

func TestDefault(t *testing.T) {
	c, err := LoadConfig(filepath.Join("testdata", "_not_exist_.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != defaultAddr {
		t.Errorf("Addr should be %q but %q", defaultAddr, c.Addr)
	}
	if c.AccessLogPath != defaultAccessLog {
		t.Errorf("AccessLogPath should be %q but %q", defaultAccessLog, c.AccessLogPath)
	}
	if c.ErrorLogPath != defaultErrorLog {
		t.Errorf("ErrorLogPath should be %q but %q", defaultErrorLog, c.ErrorLogPath)
	}
}

func TestOnlyAddr(t *testing.T) {
	c, err := LoadConfig(filepath.Join("testdata", "addr.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != "0.0.0.0:80" {
		t.Errorf("Addr should be %q but %q", "0.0.0.0:80", c.Addr)
	}
	if c.AccessLogPath != defaultAccessLog {
		t.Errorf("AccessLogPath should be %q but %q", defaultAccessLog, c.AccessLogPath)
	}
	if c.ErrorLogPath != defaultErrorLog {
		t.Errorf("ErrorLogPath should be %q but %q", defaultErrorLog, c.ErrorLogPath)
	}
}

func TestOnlyFilters(t *testing.T) {
	c, err := LoadConfig(filepath.Join("testdata", "filters.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.DefaultFilters) != 2 {
		t.Fatal("default_filters should have 2 items")
	}
	const (
		k1 = "file:///var/"
		k2 = "file:///tmp/"
		k3 = "file:///unknown/"
	)
	v1 := c.DefaultFilters[k1]
	assert.Equal(t, Filters{"tail"}, v1, "for path %q", k1)
	v2 := c.DefaultFilters[k2]
	assert.Equal(t, Filters{"head", "tail=limit:5"}, v2, "for path %q", k2)
	_, v3 := c.DefaultFilters[k3]
	if v3 {
		t.Error("Filters should be zero for unknown path")
	}
}
