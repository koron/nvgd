package core

import (
	"log"
	"os"
)

// Config represents NVGD server configuration.
type Config struct {
	Addr string
}

const defaultAddr = "127.0.0.1:9280"

// LoadConfig loads a configuration from a file.
func LoadConfig(filename string) (*Config, error) {
	if filename == "" {
		return &Config{}, nil
	}
	// TODO:
	return &Config{}, nil
}

func (c *Config) addr() string {
	if c.Addr == "" {
		return defaultAddr
	}
	return c.Addr
}

func (c *Config) logger() (*log.Logger, error) {
	// TODO: better logger.
	return log.New(os.Stderr, "", log.LstdFlags), nil
}
