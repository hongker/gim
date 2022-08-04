package config

import "gim/pkg/redis"

type Config struct {
	Redis redis.Config
	Message Message
}

func (c *Config) WithOptions(options ...Option)  {
	for _, setter := range options {
		setter(c)
	}
}

type Message struct {
	PushCount int // 每次推送的消息条数
	MaxStoreSize int // 每个会话存储的最大消息条数
}

type Option func(config *Config)

func WithMessageMaxStoreSize(maxStoreSize int) Option {
	return func(config *Config) {
		config.Message.MaxStoreSize = maxStoreSize
	}
}

func WithMessagePushCount(messagePushCount int) Option {
	return func(config *Config) {
		config.Message.PushCount = messagePushCount
	}
}

func New() *Config {
	return &Config{
		Message: Message{
			PushCount: 10,
			MaxStoreSize: 10000,
		},
	}
}