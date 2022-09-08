package generic

type Config struct {
}

func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{}
}

type CompletedConfig struct {
}

func (c CompletedConfig) New() *Server {
	return &Server{}
}
