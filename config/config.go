package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents NVGD server configuration.
type Config struct {
	Addr string `yaml:"addr"`

	// ErrorLogPath specify path of access log. default is "(stderr)".
	ErrorLogPath string `yaml:"error_log"`

	// AccessLogPath specify path of access log. default is "(discard)".
	AccessLogPath string `yaml:"access_log"`

	Protocols customConfig `yaml:"protocols"`
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
	w, err := c.openLogFile(c.AccessLogPath, "(discard)")
	if err != nil {
		return nil, err
	}
	return log.New(w, "", log.LstdFlags), nil
}

// ErrorLog creates new error logger.
func (c *Config) ErrorLog() (*log.Logger, error) {
	w, err := c.openLogFile(c.ErrorLogPath, "(stderr)")
	if err != nil {
		return nil, err
	}
	return log.New(w, "", log.LstdFlags), nil
}

func (c *Config) openLogFile(v, d string) (io.Writer, error) {
	if v == "" {
		v = d
	}
	switch v {
	case "(discard)":
		return ioutil.Discard, nil
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
	Addr:      "127.0.0.1:9280",
	Protocols: customConfig{},
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

// RegisterProtocol registers protocol configuration.
func RegisterProtocol(name string, v interface{}) {
	root.Protocols[name] = v
}
