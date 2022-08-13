package persistence

import (
	"context"
	"fmt"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisUserRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

const(
	userIdKey = "userId"
	userPrefix = "user"
)

func (repo *RedisUserRepo) getUserCacheKey(userId string) (string)  {
	return fmt.Sprintf("%s:%s", userPrefix, userId)
}

func (repo RedisUserRepo) Save(ctx context.Context, item *entity.User) error {
	err := repo.redisConn.Set(ctx, repo.getUserCacheKey(item.Id), entity.Encode(item), repo.expired).Err()
	return err
}

func (repo RedisUserRepo) Find(ctx context.Context, userId string) (*entity.User, error) {
	res, err := repo.redisConn.Get(ctx, repo.getUserCacheKey(userId)).Bytes()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}
	item := &entity.User{}
	err = entity.Decode(res, item)
	return item, err
}


func NewRedisUserRepo(redisConn redis.UniversalClient) repository.UserRepository  {
	return &RedisUserRepo{redisConn: redisConn, expired: time.Hour*24*30}
}