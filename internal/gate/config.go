package gate

import (
	"gim/pkg/grpc"
	"time"
)

type Config struct {
	TcpServer   string
	LogicServer string
	RPC         grpc.ServerConfig
}

func InitConfig() *Config {
	return &Config{
		TcpServer:   "0.0.0.0:8001",
		LogicServer: "0.0.0.0:9002",
		RPC: grpc.ServerConfig{
			Network:           "tcp",
			Addr:              "0.0.0.0:9001",
			Timeout:           time.Second * 10,
			IdleTimeout:       time.Second * 10,
			MaxLifeTime:       time.Second * 60,
			ForceCloseWait:    time.Second * 10,
			KeepAliveInterval: time.Second * 10,
			KeepAliveTimeout:  time.Second * 60,
		},
	}
}
