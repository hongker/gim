package gateway

import (
	"github.com/ebar-go/ego"
	"time"
)

type Options struct {
	ServerProtocol    string
	ServerAddress     string
	HeartbeatInterval time.Duration
}

func NewOptions() *Options {
	return &Options{
		ServerProtocol:    "tcp",
		ServerAddress:     ":8080",
		HeartbeatInterval: time.Minute,
	}
}

func (o *Options) BuildInstance() *Instance {
	c := &Config{
		HttpServerAddress: o.ServerAddress,
	}
	return &Instance{config: c, engine: ego.New()}
}
