package command

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type command struct {
	preDefined map[string]string
}

var commandHandler = &command{}

func init() {
	protocol.MustRegister("command", commandHandler)
	config.RegisterProtocol("command", &commandHandler.preDefined)
}

func (c *command) Open(u *url.URL) (*resource.Resource, error) {
	var (
		name = u.Host
	)
	cmd, ok := c.preDefined[u.Host]
	if !ok {
		return nil, fmt.Errorf("unknown command: %s", name)
	}
	rc, err := c.run(cmd)
	if err != nil {
		return nil, err
	}
	return resource.New(rc), nil
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
	return &cmdOut{stdout, cmd}, nil
}

type cmdOut struct {
	io.ReadCloser
	c *exec.Cmd
}

func (co *cmdOut) Close() error {
	return co.c.Wait()
}
