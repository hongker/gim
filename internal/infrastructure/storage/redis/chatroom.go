package redis

import (
	"context"
	"fmt"
	"gim/internal/domain/entity"
)

type ChatroomStorage struct {
	Storage
}

func (storage ChatroomStorage) cacheKey(id string) string {
	return fmt.Sprintf("chatroom:%s", id)
}

func (storage ChatroomStorage) Create(ctx context.Context, item *entity.Chatroom) error {
	bytes, _ := storage.serializer.Encode(item)
	return storage.redis.Set(ctx, storage.cacheKey(item.Id), bytes, storage.expired).Err()
}

func (storage ChatroomStorage) Find(ctx context.Context, id string) (*entity.Chatroom, error) {
	item := &entity.Chatroom{}
	bytes, err := storage.redis.Get(ctx, storage.cacheKey(id)).Bytes()
	if err != nil {
		return nil, err
	}
	err = storage.serializer.Decode(bytes, item)
	return item, err
}

func (storage ChatroomStorage) GetMember(ctx context.Context, id string) ([]string, error) {
	return storage.redis.SMembers(ctx, storage.cacheKey(id)).Result()
}

func (storage ChatroomStorage) AddMember(ctx context.Context, id string, member *entity.User) error {
	return storage.redis.SAdd(ctx, storage.cacheKey(id), member.Id).Err()
}

func (storage ChatroomStorage) RemoveMember(ctx context.Context, id string, member *entity.User) error {
	return storage.redis.SRem(ctx, storage.cacheKey(id), member.Id).Err()
}

func (storage ChatroomStorage) HasMember(ctx context.Context, id string, member *entity.User) bool {
	return storage.redis.SIsMember(ctx, storage.cacheKey(id), member.Id).Val()
}
