package gateway

import "time"

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
