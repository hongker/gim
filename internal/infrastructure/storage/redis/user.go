package redis

import (
	"context"
	"fmt"
	"gim/internal/domain/entity"
)

type UserStorage struct {
	Storage
}

func (storage UserStorage) cacheKey(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func (storage UserStorage) Create(ctx context.Context, item *entity.User) error {
	return storage.Set(ctx, storage.cacheKey(item.Id), item)
}

func (storage UserStorage) Find(ctx context.Context, id string) (*entity.User, error) {
	item := &entity.User{}
	err := storage.Get(ctx, storage.cacheKey(id), item)
	return item, err
}

func NewUserStorage() *UserStorage {
	return &UserStorage{newStorage()}
}
