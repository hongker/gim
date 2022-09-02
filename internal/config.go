package internal

type Config struct {
}

func (c *Config) Complete() *CompleteConfig {
	return &CompleteConfig{c}
}

type CompleteConfig struct {
	*Config
}

func (c *CompleteConfig) New() *Server {
	return &Server{}
}
