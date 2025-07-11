package command

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/koron/nvgd/internal/assert"
	"github.com/koron/nvgd/internal/protocoltest"
)

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "command", commandHandler)
}

func cmdRun(t *testing.T, cmd string, args ...string) string {
	t.Helper()
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		t.Fatalf("command %s %+v failed: %s", cmd, args, err)
	}
	return string(out)
}

func TestRun(t *testing.T) {
	commandHandler.preDefined = map[string]string{"goversion": "go version"}
	got := protocoltest.OpenString(t, "command://goversion")
	want := cmdRun(t, "go", "version")
	assert.Equal(t, want, got, "")
}

func TestUnknown(t *testing.T) {
	commandHandler.preDefined = nil
	got := protocoltest.OpenFail(t, "command://__unknown__")
	want := `unknown command: __unknown__`
	assert.Equal(t, want, got.Error(), "")
}

func TestEmpty(t *testing.T) {
	commandHandler.preDefined = map[string]string{"empty": ""}
	got := protocoltest.OpenFail(t, "command://empty")
	want := `empty command`
	assert.Equal(t, want, got.Error(), "")
}

func TestNotExist(t *testing.T) {
	commandHandler.preDefined = map[string]string{"notexist": "__not_exist__"}
	got := protocoltest.OpenFail(t, "command://notexist").Error()
	if !strings.HasPrefix(got, `exec: "__not_exist__": `) {
		t.Fatalf("unexpected error: %s", got)
	}
}
