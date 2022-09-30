package gateway

import "time"

type Config struct {
	Address           string
	HeartbeatInterval time.Duration
	WorkerNumber      int
	Codec             string
}

func (c *Config) New() *Controller {
	return &Controller{name: "gateway", config: c}
}
