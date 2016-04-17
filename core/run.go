package core

// Run runs NVGD server.
func Run(c *Config) error {
	s, err := New(c)
	if err != nil {
		return err
	}
	return s.Run()
}
