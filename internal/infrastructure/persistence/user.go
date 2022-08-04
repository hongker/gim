package persistence

import (
	"context"
	"fmt"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type UserRepository struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

const(
	userIdKey = "userId"
	userPrefix = "user"
)

func (repo *UserRepository) getUserCacheKey(userId string) (string)  {
	return fmt.Sprintf("%s:%s", userPrefix, userId)
}


func (repo UserRepository) Save(ctx context.Context, item *entity.User) error {
	res, err := repo.redisConn.Incr(ctx, userIdKey).Result()
	if err != nil {
		return errors.Failure(err.Error())
	}
	item.Id = strconv.FormatInt(res, 10)

	err = repo.redisConn.Set(ctx, repo.getUserCacheKey(item.Id), entity.Encode(item), repo.expired).Err()
	return nil
}

func (repo UserRepository) Get(ctx context.Context, userId string) (*entity.User, error) {
	res, err := repo.redisConn.Get(ctx, repo.getUserCacheKey(userId)).Bytes()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}
	item := &entity.User{}
	err = entity.Decode(res, item)
	return item, err
}


func NewUserRepository(redisConn redis.UniversalClient) repository.UserRepository  {
	return &UserRepository{redisConn: redisConn, expired: time.Hour*24*30}
}