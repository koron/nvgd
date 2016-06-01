package protocol

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"github.com/koron/nvgd/config"
)

type command struct {
	preDefined map[string]string
}

var commandHandler = &command{}

func init() {
	MustRegister("command", commandHandler)
	config.RegisterProtocol("command", &commandHandler.preDefined)
}

func (c *command) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		name = u.Host
	)
	cmd, ok := c.preDefined[u.Host]
	if !ok {
		return nil, fmt.Errorf("unknown command: %s", name)
	}
	return c.run(cmd)
}

func (c *command) run(s string) (io.ReadCloser, error) {
	ss := strings.Split(s, " ")
	if len(s) < 1 {
		return nil, errors.New("empty command")
	}
	cmd := exec.Command(ss[0], ss[1:]...)
	// clone STDOUT
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	// FIXME: clone STDERR and merge to output
	// start the command.
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, err
	}
	return stdout, nil
}
