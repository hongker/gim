package internal

import (
	"gim/internal/controller/api"
	"gim/internal/controller/gateway"
	"gim/internal/controller/job"
)

type Config struct {
	// ApiControllerConfig
	ApiControllerConfig *api.Config

	// GatewayControllerConfig
	GatewayControllerConfig *gateway.Config

	// JobControllerConfig
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
func (config *Config) BuildInstance() *Server {
	return &Server{
		config: config,
	}
}
