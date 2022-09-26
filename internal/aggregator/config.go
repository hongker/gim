package aggregator

import (
	"gim/internal/controllers/api"
	"gim/internal/controllers/gateway"
)

type Config struct {
	// ApiControllerConfig
	ApiControllerConfig *api.Config

	// GatewayControllerConfig
	GatewayControllerConfig *gateway.Config
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	return &Config{
		ApiControllerConfig:     &api.Config{},
		GatewayControllerConfig: &gateway.Config{},
	}
}

// New build aggregator instance.
func (c *Config) New() *Aggregator {
	return &Aggregator{
		config: c,
	}
}
