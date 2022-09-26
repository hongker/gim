package aggregator

import (
	"gim/internal/controllers/api"
	"gim/internal/controllers/gateway"
)

type Config struct {
	ApiControllerConfig *api.Config

	GatewayControllerConfig *gateway.Config
}

func NewConfig() *Config {
	return &Config{
		ApiControllerConfig:     &api.Config{},
		GatewayControllerConfig: &gateway.Config{},
	}
}
func (c *Config) New() *Aggregator {
	return &Aggregator{
		config: c,
	}
}
