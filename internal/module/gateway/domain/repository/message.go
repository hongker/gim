package repository

import (
	"context"
	"gim/internal/module/gateway/domain/entity"
	"gim/internal/module/gateway/domain/types"
	"github.com/ebar-go/ego/errors"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *types.Message) error
	Find(ctx context.Context, id string) (*entity.Message, error)
}

func NewMessageRepository() MessageRepository {
	return &messageRepo{
		messages: make(map[string]*entity.Message),
	}
}

type messageRepo struct {
	mu       sync.RWMutex
	messages map[string]*entity.Message
}

func (repo *messageRepo) Save(ctx context.Context, msg *types.Message) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	msg.Id = uuid.NewV4().String()
	repo.messages[msg.Id] = &entity.Message{
		Id:          msg.Id,
		SenderId:    msg.SenderId,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		Status:      msg.Status,
		CreatedAt:   msg.CreatedAt,
	}
	return nil
}
func (repo *messageRepo) Find(ctx context.Context, id string) (*entity.Message, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	item, ok := repo.messages[id]
	if !ok {
		return nil, errors.NotFound("message not found")
	}
	return item, nil
}
