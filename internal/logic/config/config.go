package config

import (
	"gim/pkg/grpc"
	"time"
)

type Config struct {
	RPC grpc.ServerConfig
}

func Init() *Config {
	return &Config{
		RPC: grpc.ServerConfig{
			Network:           "tcp",
			Addr:              "0.0.0.0:9002",
			Timeout:           time.Second * 10,
			IdleTimeout:       time.Second * 10,
			MaxLifeTime:       time.Second * 60,
			ForceCloseWait:    time.Second * 10,
			KeepAliveInterval: time.Second * 10,
			KeepAliveTimeout:  time.Second * 60,
		},
	}
}
