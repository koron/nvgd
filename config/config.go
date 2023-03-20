// Package config provides configuration object for NVGD.
package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents NVGD server configuration.
type Config struct {
	Addr string `yaml:"addr"`

	PathPrefix string `yaml:"path_prefix"`

	// ErrorLogPath specify path of access log. default is "(stderr)".
	ErrorLogPath string `yaml:"error_log"`

	// AccessLogPath specify path of access log. default is "(discard)".
	AccessLogPath string `yaml:"access_log"`

	Protocols customConfig `yaml:"protocols,omitempty"`

	Filters customConfig `yaml:"filters,omitempty"`

	DefaultFilters FiltersMap `yaml:"default_filters,omitempty"`

	// Aliases provides custom aliases.
	Aliases map[string]string `yaml:"aliases,omitempty"`

	// AccessControlAllowOrigin specify a value for
	// "Access-Control-Allow-Origin" header put to all responses.
	// The header is not put when this is empty string.
	// Default value is empty.
	AccessControlAllowOrigin string `yaml:"access_control_allow_origin,omitempty"`
}

type customConfig map[string]interface{}

func (cc customConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m yaml.MapSlice
	if err := unmarshal(&m); err != nil {
		return err
	}
	for _, item := range m {
		k, ok := item.Key.(string)
		if !ok {
			continue
		}
		v, ok := cc[k]
		if !ok {
			return fmt.Errorf("unknown configuration name: %s", k)
		}
		b, err := yaml.Marshal(item.Value)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, v); err != nil {
			return err
		}
	}
	return nil
}

// AccessLog creates a new access logger.
func (c *Config) AccessLog() (*log.Logger, error) {
	w, err := c.openLogFile(c.AccessLogPath)
	if err != nil {
		return nil, err
	}
	return log.New(w, "", log.LstdFlags), nil
}

// ErrorLog creates new error logger.
func (c *Config) ErrorLog() (*log.Logger, error) {
	w, err := c.openLogFile(c.ErrorLogPath)
	if err != nil {
		return nil, err
	}
	return log.New(w, "", log.LstdFlags), nil
}

func (c *Config) openLogFile(v string) (io.Writer, error) {
	switch v {
	case "(discard)":
		return io.Discard, nil
	case "(stderr)":
		return os.Stderr, nil
	case "(stdout)":
		return os.Stdout, nil
	default:
		f, err := os.OpenFile(v, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

var root = &Config{
	Addr:           defaultAddr,
	AccessLogPath:  defaultAccessLog,
	ErrorLogPath:   defaultErrorLog,
	Protocols:      customConfig{},
	Filters:        customConfig{},
	DefaultFilters: FiltersMap{},
}

func Root() *Config {
	return root
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
	b, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, root); err != nil {
		return nil, err
	}
	return root, nil
}

// RegisterProtocol registers protocol configuration.
func RegisterProtocol(name string, v interface{}) {
	root.Protocols[name] = v
}

// RegisterFilter registers a filter configuration.
func RegisterFilter(name string, v interface{}) {
	root.Filters[name] = v
}
