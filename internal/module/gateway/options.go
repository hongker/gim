package gateway

import (
	"github.com/ebar-go/ego"
	"time"
)

type Options struct {
	HttpServerAddress string
	GrpcServerAddress string
	SockServerAddress string
	HeartbeatInterval time.Duration
}

func NewOptions() *Options {
	return &Options{
		HttpServerAddress: ":8080",
		GrpcServerAddress: ":8081",
		SockServerAddress: ":8082",
		HeartbeatInterval: time.Minute,
	}
}

func (o *Options) BuildInstance() *Instance {
	c := &Config{
		HttpServerAddress: o.HttpServerAddress,
		GrpcServerAddress: o.GrpcServerAddress,
		SockServerAddress: o.SockServerAddress,
	}
	return &Instance{config: c, engine: ego.New()}
}
