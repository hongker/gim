package repository

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/types"
	"gim/internal/infra/storage"
	uuid "github.com/satori/go.uuid"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *types.Message) error
	Find(ctx context.Context, id string) (*entity.Message, error)
}

func NewMessageRepository() MessageRepository {
	return &messageRepo{
		store: storage.NewMemoryStorage("message"),
	}
}

type messageRepo struct {
	store storage.Storage
}

func (repo *messageRepo) Save(ctx context.Context, msg *types.Message) error {
	item := &entity.Message{
		SenderId:  msg.SenderId,
		Content:   msg.Content,
		Category:  string(msg.Category),
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
	}
	item.Id = uuid.NewV4().String()

	if err := repo.store.Save(ctx, item); err != nil {
		return err
	}
	msg.Id = item.ID()
	return nil
}

func (repo *messageRepo) Find(ctx context.Context, id string) (*entity.Message, error) {
	item := entity.NewMessageWithID(id)
	err := repo.store.Find(ctx, item)
	return item, err
}
