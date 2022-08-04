package config

import (
	"fmt"
	"gim/pkg/redis"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	viper *viper.Viper
	Server  Server
	Redis redis.Config
	Message Message
}

func (c *Config) LoadFile(path ...string) (err error) {
	for _, p := range path {
		c.viper.SetConfigFile(p)
		if err = c.viper.MergeInConfig(); err != nil {
			return
		}
	}

	return c.viper.Unmarshal(c)
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
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

type Server struct {
	Protocol string //
	Port int //
}


type Option func(config *Config)

func New() *Config {
	return &Config{
		viper: viper.New(),
		Server: Server{
			Protocol: "tcp",
			Port:     8080,
		},
		Redis: redis.Config{
			Host:        "127.0.0.1",
			Port:        6379,
			Auth:        "",
			PoolSize:    10,
			MaxRetries:  3,
			IdleTimeout: time.Second * 10,
			Cluster:     nil,
		},
		Message: Message{
			PushCount: 10,
			MaxStoreSize: 10000,
		},
	}
}