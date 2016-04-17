package core

// Config represents NVGD server configuration.
type Config struct {
	Addr string
}

const default_addr = "127.0.0.1:9280"

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
		return default_addr
	}
	return c.Addr
}
