package options

import (
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
)

type ServerRunOptions struct {
	Protocol            string
	Port                int
	Debug               bool
	MessageMaxStoreSize int
	MessagePushCount    int
	MessageStorage      string
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		Protocol:            "tcp",
		Port:                8080,
		Debug:               false,
		MessageMaxStoreSize: 1000,
		MessagePushCount:    5,
		MessageStorage:      infrastructure.MemoryStore,
	}
	return s
}

func (s ServerRunOptions) ApplyTo(conf *config.Config) {
	conf.Server.Protocol = s.Protocol
	conf.Server.Port = s.Port
	conf.Server.Store = s.MessageStorage
	conf.Message.PushCount = s.MessagePushCount
	conf.Message.MaxStoreSize = s.MessageMaxStoreSize
}
