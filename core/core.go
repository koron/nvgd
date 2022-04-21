// Package core provides core logic for NVGD.
package core

import "github.com/koron/nvgd/config"

// Run runs NVGD server.
func Run(c *config.Config) error {
	s, err := New(c)
	if err != nil {
		return err
	}
	return s.Run()
}
