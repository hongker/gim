package memory

import (
	"context"
	"gim/internal/domain/entity"
	"sync"
)

type MappingStorage[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func (storage *MappingStorage[T]) Get(key string) (T, bool) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	item, ok := storage.items[key]
	return item, ok
}

type UserStorage struct {
}

func (storage *UserStorage) Create(ctx context.Context, item *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (storage *UserStorage) Find(ctx context.Context, id string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
