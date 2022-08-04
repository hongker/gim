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

type GroupRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

const(
	groupCachePrefix = "group"
	groupUserCachePrefix = "groupUser"
	groupIdKey = "groupId"

)

func (repo GroupRepo) getCacheKey(groupId string) string {
	return fmt.Sprintf("%s:%s", groupCachePrefix, groupId)
}

func (repo GroupRepo) Find(ctx context.Context, id string) (*entity.Group, error) {
	res, err := repo.redisConn.Get(ctx, repo.getCacheKey(id)).Bytes()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}
	item := &entity.Group{}
	err = entity.Decode(res, item)
	return item, err
}

func (repo GroupRepo) Create(ctx context.Context, item *entity.Group) error {
	res, err := repo.redisConn.Incr(ctx, groupIdKey).Result()
	if err != nil {
		return errors.Failure(err.Error())
	}
	item.Id = strconv.FormatInt(res, 10)

	err = repo.redisConn.Set(ctx, repo.getCacheKey(item.Id), entity.Encode(item), repo.expired).Err()
	return nil
}

func NewGroupRepo(redisConn redis.UniversalClient) repository.GroupRepo  {
	return &GroupRepo{redisConn: redisConn, expired: time.Hour* 24 * 30}
}

type GroupUserRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

func (repo GroupUserRepo) getCacheKey(groupId, userId string) string {
	return fmt.Sprintf("%s:%s:%s", groupUserCachePrefix, groupId, userId)
}

func (repo GroupUserRepo) Create(ctx context.Context, item *entity.GroupUser) error {
	err := repo.redisConn.Set(ctx,  repo.getCacheKey(item.GroupId, item.UserId), entity.Encode(item), repo.expired).Err()
	return err
}

func (repo GroupUserRepo) Find(ctx context.Context, groupId string, userId string) (*entity.GroupUser, error) {
	res, err := repo.redisConn.Get(ctx, repo.getCacheKey(groupId, userId)).Bytes()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}
	item := &entity.GroupUser{}
	err = entity.Decode(res, item)
	return item, err
}

func (repo GroupUserRepo) Delete(ctx context.Context, groupId string, userId string) error {
	return repo.redisConn.Del(ctx, repo.getCacheKey(groupId, userId)).Err()
}

func NewGroupUserRepo(redisConn redis.UniversalClient) repository.GroupUserRepo {
	return &GroupUserRepo{redisConn: redisConn, expired: time.Hour* 24 * 30}
}
