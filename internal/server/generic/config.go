package generic

import (
	"gim/internal/application"
	"gim/pkg/network"
)

type Config struct {
	TCPAddr string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{
		Config: c,
	}
}

type CompletedConfig struct {
	*Config
}

func (c CompletedConfig) New() *Server {
	socketInstance := network.NewTCPServer(c.TCPAddr)

	return &Server{
		el: application.BuildEventLoop(socketInstance),
	}
}
