package repository

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/infrastructure/storage"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
)

type SessionRepository interface {
	List(ctx context.Context, uid string) ([]*entity.Session, error)
	SaveMessage(ctx context.Context, uid string, session *entity.Session, msg *entity.Message) error
	QueryMessage(ctx context.Context, session *entity.Session) ([]*entity.Message, error)
	FindMessage(ctx context.Context, id string) (*entity.Message, error)
}

type sessionRepo struct {
	store *storage.StorageManager
}

func (s *sessionRepo) FindMessage(ctx context.Context, id string) (*entity.Message, error) {
	return s.store.Message().Find(ctx, id)
}

func (repo *sessionRepo) List(ctx context.Context, uid string) ([]*entity.Session, error) {
	items, err := repo.store.Session().List(ctx, uid)
	return items, err
}

func (repo *sessionRepo) QueryMessage(ctx context.Context, session *entity.Session) ([]*entity.Message, error) {
	ids, err := repo.store.Session().QueryMessage(ctx, session.Id)
	if err != nil {
		return nil, err
	}
	messages := make([]*entity.Message, 0, len(ids))
	for _, id := range ids {
		item, lastErr := repo.store.Message().Find(ctx, id)
		if lastErr != nil {
			continue
		}
		messages = append(messages, item)
	}
	return messages, nil
}

func (repo *sessionRepo) SaveMessage(ctx context.Context, uid string, session *entity.Session, msg *entity.Message) error {
	return runtime.Call(func() error {
		err := repo.store.Session().Create(ctx, uid, session)
		if err == nil {
			return nil
		}
		return errors.WithMessage(err, "create user session")
	}, func() error {
		err := repo.store.Message().Create(ctx, msg)
		if err == nil {
			return nil
		}
		return errors.WithMessage(err, "create message")
	}, func() error {
		err := repo.store.Session().SaveMessage(ctx, session.Id, msg.Id)
		if err == nil {
			return nil
		}
		return errors.WithMessage(err, "save session message")
	})

}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{store: storage.MemoryManager()}
}
