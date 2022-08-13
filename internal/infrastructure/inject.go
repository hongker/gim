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


const(
	RedisStore = "redis"
	MemoryStore = "memory"
	MongoStore = "mongodb"
)
func newRedis(conf *config.Config) (redis.UniversalClient, error) {
	return gredis.Connect(conf.Redis)
}

func InjectStore(container *dig.Container, store string)  {
	_ = container.Provide(config.New)
	_ = container.Provide(cache.NewFactory)

	if store == MemoryStore { // 默认使用内存缓存
		_ = container.Provide(func( factory *cache.Factory) repository.UserRepository{
			return cache.NewUserRepo(nil, factory.Create())
		})
		_ = container.Provide(func( factory *cache.Factory) repository.GroupRepo{
			return cache.NewGroupRepo(nil, factory.Create())
		})
		_ = container.Provide(func(factory *cache.Factory) repository.GroupUserRepo{
			return cache.NewGroupUserRepo(nil)
		})
		_ = container.Provide(func(factory *cache.Factory) repository.MessageRepo{
			return cache.NewMessageRepo(nil)
		})
	}else if store == RedisStore {
		_ = container.Provide(newRedis)

		_ = container.Provide(func(redisConn redis.UniversalClient,  factory *cache.Factory) repository.UserRepository{
			delegate := persistence.NewRedisUserRepo(redisConn)
			return cache.NewUserRepo(delegate, factory.Create())
		})

		_ = container.Provide(func(redisConn redis.UniversalClient,  factory *cache.Factory) repository.GroupRepo{
			delegate := persistence.NewRedisGroupRepo(redisConn)
			return cache.NewGroupRepo(delegate, factory.Create())
		})

		_ = container.Provide(func(redisConn redis.UniversalClient,  factory *cache.Factory) repository.GroupUserRepo{
			delegate := persistence.NewRedisGroupUserRepo(redisConn)
			return cache.NewGroupUserRepo(delegate)
		})
		_ = container.Provide(func(redisConn redis.UniversalClient,  factory *cache.Factory) repository.MessageRepo{
			delegate := persistence.NewRedisMessageRepo(redisConn)
			return delegate
		})
	}else if store == MongoStore {

	}
}
