package infrastructure

import (
	"gim/internal/domain/repository"
	"gim/internal/infrastructure/cache"
	"gim/internal/infrastructure/config"
	"gim/internal/infrastructure/persistence"
	gredis "gim/pkg/redis"
	"github.com/go-redis/redis/v8"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(config.New)
	_ = container.Provide(cache.NewFactory)
	_ = container.Provide(newRedis)
	_ = container.Provide(persistence.NewMessageRepo)
	_ = container.Provide(newUserRepository)
	_ = container.Provide(newGroupRepository)
	_ = container.Provide(persistence.NewGroupUserRepo)
}

func newRedis(conf *config.Config) (redis.UniversalClient, error) {
	return gredis.Connect(conf.Redis)
}

func newGroupRepository(redisConn redis.UniversalClient, factory cache.Factory) repository.GroupRepo   {
	delegate := persistence.NewGroupRepo(redisConn)
	return cache.NewGroupRepo(delegate, factory.Create())
}


func newUserRepository(redisConn redis.UniversalClient, factory cache.Factory) repository.UserRepository   {
	delegate := persistence.NewUserRepository(redisConn)
	return cache.NewUserRepo(delegate, factory.Create())
}

