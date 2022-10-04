package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"gim/internal/domain/entity"
	"github.com/ebar-go/ego/component"
	"time"
)

type Storage struct {
	expired    time.Duration
	redis      *component.Redis
	serializer *Serializer
}

func (storage Storage) Set(ctx context.Context, key string, value any) error {
	bytes, err := storage.serializer.Encode(value)
	if err != nil {
		return err
	}
	return storage.redis.Set(ctx, key, bytes, storage.expired).Err()
}
func (storage Storage) Get(ctx context.Context, key string, value any) error {
	bytes, err := storage.redis.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return storage.serializer.Decode(bytes, value)
}

type Serializer struct {
}

func (s *Serializer) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (s *Serializer) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

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
