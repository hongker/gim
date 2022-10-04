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
	encoded, _ := storage.serializer.Encode(item)
	return storage.redis.Set(ctx, storage.cacheKey(item.Id), encoded, storage.expired).Err()
}

func (storage UserStorage) Find(ctx context.Context, id string) (*entity.User, error) {
	bytes, err := storage.redis.Get(ctx, id).Bytes()
	if err != nil {
		return nil, err
	}
	item := &entity.User{}
	err = storage.serializer.Decode(bytes, item)
	return item, err
}
