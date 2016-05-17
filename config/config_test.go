package config

import "testing"

func TestDefault(t *testing.T) {
	c, err := LoadConfig("empty.yml")
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
	c, err := LoadConfig("addr.yml")
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
