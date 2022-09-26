package internal

import (
	"gim/internal/controllers/api"
	"gim/internal/controllers/gateway"
	"gim/internal/controllers/job"
)

type Config struct {
	// ApiControllerConfig
	ApiControllerConfig *api.Config

	// GatewayControllerConfig
	GatewayControllerConfig *gateway.Config

	JobControllerConfig *job.Config
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	return &Config{
		ApiControllerConfig:     &api.Config{},
		GatewayControllerConfig: &gateway.Config{},
		JobControllerConfig:     &job.Config{},
	}
}

// BuildInstance return a new server instance.
func (c *Config) BuildInstance() *Server {
	return &Server{
		config: c,
	}
}
