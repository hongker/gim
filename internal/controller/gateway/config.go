package gateway

import "time"

type Config struct {
	Address           string
	HeartbeatInterval time.Duration
	WorkerNumber      int
	Codec             string
}

func (c *Config) New(name string) *Controller {
	return &Controller{name: name, config: c}
}
