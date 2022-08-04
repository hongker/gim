package infrastructure

import (
	"gim/internal/infrastructure/config"
	"gim/internal/infrastructure/persistence"
	"gim/pkg/redis"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(config.New)
	_ = container.Provide(func(conf *config.Config) redis.Config{
		return conf.Redis
	})
	_ = container.Provide(redis.Connect)
	_ = container.Provide(persistence.NewMessageRepo)
	_ = container.Provide(persistence.NewUserRepository)
	_ = container.Provide(persistence.NewGroupRepo)
	_ = container.Provide(persistence.NewGroupUserRepo)
}
