package storage

import (
	"context"
	"github.com/ebar-go/ego/errors"
	"sync"
)

func NewMemoryStorage(name string) Storage {
	return &MemoryStorage{name: name, items: map[string]Object{}}
}

type MemoryStorage struct {
	name  string
	mu    sync.RWMutex // guards
	items map[string]Object
}

func (storage *MemoryStorage) Save(ctx context.Context, object Object) error {
	storage.mu.Lock()
	storage.items[object.ID()] = object
	storage.mu.Unlock()
	return nil
}
func (storage *MemoryStorage) Find(ctx context.Context, object Object) error {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	var found bool
	object, found = storage.items[object.ID()]
	if !found {
		return errors.NotFound("object not found")
	}

	return nil
}
