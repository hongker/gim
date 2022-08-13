package infrastructure

import (
	"gim/internal/domain/repository"
	"gim/internal/infrastructure/cache"
	"gim/internal/infrastructure/config"
	"gim/internal/infrastructure/persistence"
	"gim/pkg/redis"
	goredis "github.com/go-redis/redis/v8"
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
	_ = container.Provide(newGroupRepository)
	_ = container.Provide(persistence.NewGroupUserRepo)
}

func newGroupRepository(redisConn goredis.UniversalClient, conf *config.Config) repository.GroupRepo   {
	delegate := persistence.NewGroupRepo(redisConn)
	return cache.NewGroupRepo(delegate, conf)
}


func newUserRepository(redisConn goredis.UniversalClient, conf *config.Config) repository.UserRepository   {
	delegate := persistence.NewUserRepository(redisConn)
	return cache.NewUserRepo(delegate, conf)
}

