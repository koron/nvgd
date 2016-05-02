package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents NVGD server configuration.
type Config struct {
	Addr string
}

var root = &Config{
	Addr: "127.0.0.1:9280",
}

// LoadConfig loads a configuration from a file.
func LoadConfig(filename string) (*Config, error) {
	if filename == "" {
		return root, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return root, nil
		}
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, root); err != nil {
		return nil, err
	}
	return root, nil
}

// GetLogger gets logger.
func (c *Config) GetLogger() (*log.Logger, error) {
	// TODO: better logger.
	return log.New(os.Stderr, "", log.LstdFlags), nil
}
