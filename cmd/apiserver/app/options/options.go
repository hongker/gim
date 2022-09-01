package options

import (
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/pkg/redis"
	"time"
)

type ServerRunOptions struct {
	Protocol            string
	Port                int
	Debug               bool
	MessageMaxStoreSize int
	MessagePushCount    int
	MessageStorage      string
	HeartbeatInterval   time.Duration
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		Protocol:            "tcp",
		Port:                8080,
		Debug:               false,
		MessageMaxStoreSize: 1000,
		MessagePushCount:    5,
		MessageStorage:      infrastructure.MemoryStore,
		HeartbeatInterval:   time.Minute,
	}
	return s
}

func (s ServerRunOptions) ApplyTo(conf *config.Config) {
	conf.Debug = s.Debug
	conf.Server = config.Server{
		Protocol:          s.Protocol,
		Port:              s.Port,
		HeartbeatInterval: s.HeartbeatInterval,
		Store:             s.MessageStorage,
	}
	conf.Redis = redis.Config{
		Host:        "127.0.0.1",
		Port:        6379,
		Auth:        "",
		PoolSize:    10,
		MaxRetries:  3,
		IdleTimeout: time.Second * 10,
		Cluster:     nil,
	}
	conf.Message = config.Message{
		PushCount:    s.MessagePushCount,
		MaxStoreSize: s.MessageMaxStoreSize,
	}
}
