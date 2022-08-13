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
	_ = container.Provide(newRedis)
	_ = container.Provide(persistence.NewMessageRepo)
	_ = container.Provide(newUserRepository)
	_ = container.Provide(newGroupRepository)
	_ = container.Provide(persistence.NewGroupUserRepo)
}

func newRedis(conf *config.Config) (goredis.UniversalClient, error) {
	return redis.Connect(conf.Redis)
}

func newGroupRepository(redisConn goredis.UniversalClient, conf *config.Config) repository.GroupRepo   {
	delegate := persistence.NewGroupRepo(redisConn)
	return cache.NewGroupRepo(delegate, conf)
}


func newUserRepository(redisConn goredis.UniversalClient, conf *config.Config) repository.UserRepository   {
	delegate := persistence.NewUserRepository(redisConn)
	return cache.NewUserRepo(delegate, conf)
}

