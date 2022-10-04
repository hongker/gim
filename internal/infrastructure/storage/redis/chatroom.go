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
	return storage.Set(ctx, storage.cacheKey(item.Id), item)
}

func (storage ChatroomStorage) Find(ctx context.Context, id string) (*entity.Chatroom, error) {
	item := &entity.Chatroom{}
	err := storage.Get(ctx, storage.cacheKey(id), item)
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

func NewChatroomStorage() *ChatroomStorage {
	return &ChatroomStorage{newStorage()}
}
