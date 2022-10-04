package redis

import (
	"context"
	"fmt"
	"gim/internal/domain/entity"
)

type SessionStorage struct {
	Storage
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{newStorage()}
}

func (storage SessionStorage) cacheKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (storage SessionStorage) sessionMessageCacheKey(id string) string {
	return fmt.Sprintf("session:message:%s", id)
}

func (storage SessionStorage) Create(ctx context.Context, uid string, item *entity.Session) error {
	bytes, _ := storage.serializer.Encode(item)
	return storage.redis.HSet(ctx, storage.cacheKey(uid), item.Id, bytes).Err()
}

func (storage SessionStorage) List(ctx context.Context, uid string) ([]*entity.Session, error) {
	items := storage.redis.HGetAll(ctx, storage.cacheKey(uid)).Val()

	res := make([]*entity.Session, 0)
	if len(items) == 0 {
		return res, nil
	}

	for _, item := range items {
		session := &entity.Session{}
		if err := storage.serializer.Decode([]byte(item), session); err != nil {
			continue
		}
		res = append(res, session)
	}
	return res, nil
}

func (storage SessionStorage) QueryMessage(ctx context.Context, sessionId string) ([]string, error) {
	return storage.redis.LRange(ctx, storage.sessionMessageCacheKey(sessionId), 0, -1).Result()
}

func (storage SessionStorage) SaveMessage(ctx context.Context, sessionId, msgId string) error {
	return storage.redis.LPush(ctx, storage.sessionMessageCacheKey(sessionId), msgId).Err()
}

type MessageStorage struct {
	Storage
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{newStorage()}
}

func (storage *MessageStorage) cacheKey(id string) string {
	return fmt.Sprintf("message:%s", id)
}

func (storage MessageStorage) Create(ctx context.Context, msg *entity.Message) error {
	return storage.Set(ctx, storage.cacheKey(msg.Id), msg)
}

func (storage MessageStorage) Find(ctx context.Context, id string) (*entity.Message, error) {
	item := &entity.Message{}
	err := storage.Get(ctx, storage.cacheKey(id), item)
	return item, err
}
