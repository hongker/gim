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

type RedisGroupRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

const(
	groupCachePrefix = "group"
	groupUserCachePrefix = "groupUser"
	groupIdKey = "groupId"

)

func (repo RedisGroupRepo) getCacheKey(groupId string) string {
	return fmt.Sprintf("%s:%s", groupCachePrefix, groupId)
}

func (repo RedisGroupRepo) Find(ctx context.Context, id string) (*entity.Group, error) {
	res, err := repo.redisConn.Get(ctx, repo.getCacheKey(id)).Bytes()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}
	item := &entity.Group{}
	err = entity.Decode(res, item)
	return item, err
}

func (repo RedisGroupRepo) Create(ctx context.Context, item *entity.Group) error {
	res, err := repo.redisConn.Incr(ctx, groupIdKey).Result()
	if err != nil {
		return errors.Failure(err.Error())
	}
	item.Id = strconv.FormatInt(res, 10)

	err = repo.redisConn.Set(ctx, repo.getCacheKey(item.Id), entity.Encode(item), repo.expired).Err()
	return nil
}

func NewRedisGroupRepo(redisConn redis.UniversalClient) repository.GroupRepo  {
	return &RedisGroupRepo{redisConn: redisConn, expired: time.Hour* 24 * 30}
}

type RedisGroupUserRepo struct {
	redisConn redis.UniversalClient
	expired time.Duration
}

func (repo RedisGroupUserRepo) getCacheKey(groupId string) string {
	return fmt.Sprintf("%s:%s", groupUserCachePrefix, groupId)
}

func (repo RedisGroupUserRepo) FindAll(ctx context.Context, groupId string) ([]string, error){
	res, err := repo.redisConn.SMembers(ctx, repo.getCacheKey(groupId)).Result()
	return res, err
}

func (repo RedisGroupUserRepo) Create(ctx context.Context, item *entity.GroupUser) error {
	err := repo.redisConn.SAdd(ctx,  repo.getCacheKey(item.GroupId), item.UserId).Err()
	return err
}

func (repo RedisGroupUserRepo) Find(ctx context.Context, groupId string, userId string) (*entity.GroupUser, error) {
	res, err := repo.redisConn.SIsMember(ctx, repo.getCacheKey(groupId), userId).Result()
	if err != nil {
		return nil, errors.Failure(err.Error())
	}

	if !res {
		return nil, errors.DataNotFound("user not found")
	}
	item := &entity.GroupUser{
		GroupId:   groupId,
		UserId:    userId,
		CreatedAt: 0,
	}
	return item, err
}

func (repo RedisGroupUserRepo) Delete(ctx context.Context, groupId string, userId string) error {
	return repo.redisConn.SRem(ctx, repo.getCacheKey(groupId), userId).Err()
}

func NewRedisGroupUserRepo(redisConn redis.UniversalClient) repository.GroupUserRepo {
	return &RedisGroupUserRepo{redisConn: redisConn, expired: time.Hour* 24 * 30}
}
