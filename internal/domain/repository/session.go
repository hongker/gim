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
	QueryMessage(ctx context.Context, session *entity.Session)
}

type sessionRepo struct {
	store *storage.StorageManager
}

func (repo *sessionRepo) List(ctx context.Context, uid string) ([]*entity.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *sessionRepo) QueryMessage(ctx context.Context, session *entity.Session) {
	//TODO implement me
	panic("implement me")
}

func (repo *sessionRepo) SaveMessage(ctx context.Context, uid string, session *entity.Session, msg *entity.Message) error {
	return runtime.Call(func() error {
		err := repo.store.Session().Create(ctx, uid, session)
		return errors.WithMessage(err, "create user session")
	}, func() error {
		err := repo.store.Message().Create(ctx, msg)
		return errors.WithMessage(err, "create message")
	}, func() error {
		err := repo.store.Session().SaveMessage(ctx, session.Id, msg.Id)
		return errors.WithMessage(err, "save session message")
	})

}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{store: storage.MemoryManager()}
}
