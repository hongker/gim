package internal

type Config struct {
	generic GenericServerConfig
	message MessageConfig
}

func New() *Config {
	return &Config{}
}

type GenericServerConfig struct {
}

type MessageConfig struct {
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
