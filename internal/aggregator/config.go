package aggregator

import (
	"gim/internal/controllers/api"
	"gim/internal/controllers/socket"
)

type Config struct {
	apiControllerConfig *api.Config

	gatewayControllerConfig *socket.Config
}

func NewConfig() *Config {
	return &Config{
		apiControllerConfig: &api.Config{
			Address:         ":8080",
			TraceHeader:     "trace",
			EnableProfiling: false,
		},
		gatewayControllerConfig: &socket.Config{
			Address: ":8081",
		},
	}
}
func (c *Config) New() *Aggregator {
	return &Aggregator{
		config: c,
	}
}
