package redis

import (
	"context"
	"encoding/json"
	"github.com/ebar-go/ego/component"
	"time"
)

var (
	defaultSerializer = &Serializer{}
	defaultExpired    = time.Hour * 24 * 30
)

type Storage struct {
	expired    time.Duration
	redis      *component.Redis
	serializer *Serializer
}

func newStorage() Storage {
	return Storage{
		expired:    defaultExpired,
		redis:      component.Provider().Redis(),
		serializer: defaultSerializer,
	}
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
