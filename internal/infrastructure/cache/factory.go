package cache

import (
	"gim/internal/infrastructure/config"
	"github.com/patrickmn/go-cache"
	"time"
)

type Factory struct {
	defaultExpiration time.Duration
	cleanupInterval time.Duration
}

func NewFactory(conf *config.Config) *Factory {
	return &Factory{
		defaultExpiration: conf.Cache.Expired,
		cleanupInterval:   conf.Cache.Cleanup,
	}
}

func (f *Factory) Create() *cache.Cache  {
	return cache.New(f.defaultExpiration, f.cleanupInterval)
}
