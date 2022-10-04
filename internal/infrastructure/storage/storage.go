package storage

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/infrastructure/storage/memory"
	"gim/internal/infrastructure/storage/redis"
	"sync"
)

type Message interface {
	Create(ctx context.Context, msg *entity.Message) error
	Find(ctx context.Context, id string) (*entity.Message, error)
}
type User interface {
	Create(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}
type Chatroom interface {
	Create(ctx context.Context, item *entity.Chatroom) error
	Find(ctx context.Context, id string) (*entity.Chatroom, error)
	GetMember(ctx context.Context, id string) ([]string, error)
	AddMember(ctx context.Context, id string, member *entity.User) error
	RemoveMember(ctx context.Context, id string, member *entity.User) error
	HasMember(ctx context.Context, id string, member *entity.User) bool
}

type Session interface {
	Create(ctx context.Context, uid string, item *entity.Session) error
	List(ctx context.Context, uid string) ([]*entity.Session, error)
	QueryMessage(ctx context.Context, sessionId string) ([]string, error)
	SaveMessage(ctx context.Context, sessionId, msgId string) error
}

type StorageManager struct {
	user     User
	message  Message
	chatroom Chatroom
	session  Session
}

func (s StorageManager) User() User {
	return s.user
}

func (s StorageManager) Message() Message {
	return s.message
}

func (s StorageManager) Chatroom() Chatroom {
	return s.chatroom
}

func (s StorageManager) Session() Session {
	return s.session
}

var storageProvider = struct {
	once     sync.Once
	instance *StorageManager
}{}

func MemoryManager() *StorageManager {
	storageProvider.once.Do(func() {
		storageProvider.instance = &StorageManager{
			user:     memory.NewUserStorage(),
			message:  memory.NewMessageStorage(),
			chatroom: memory.NewChatroomStorage(),
			session:  memory.NewSessionStorage(),
		}
	})
	return storageProvider.instance

}

func RedisManager() *StorageManager {
	storageProvider.once.Do(func() {
		storageProvider.instance = &StorageManager{
			user:     redis.NewUserStorage(),
			message:  redis.NewMessageStorage(),
			chatroom: redis.NewChatroomStorage(),
			session:  redis.NewSessionStorage(),
		}
	})
	return storageProvider.instance
}
