package repository

import (
	"context"
	types "gim/internal/domain/types"
	"gim/internal/infra/storage"
	"github.com/ebar-go/ego/errors"
)

type SessionRepository interface {
	List(ctx context.Context, uid string) ([]types.Session, error)
	SaveMessage(ctx context.Context, session *types.Session, msg *types.Message) error
	QueryMessage(ctx context.Context, session *types.Session)
}

type sessionRepo struct {
	store storage.Storage
}

func (repo *sessionRepo) List(ctx context.Context, uid string) ([]types.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *sessionRepo) SaveMessage(ctx context.Context, session *types.Session, msg *types.Message) error {
	sessionMessages := &types.SessionMessage{Id: session.Id}
	if err := repo.store.Find(ctx, sessionMessages); err != nil {
		if !errors.Is(err, errors.NotFound("")) {
			return err
		}
	}
	sessionMessages.AddMessage(msg.Id)
	return repo.store.Save(ctx, sessionMessages)
}

func (repo *sessionRepo) QueryMessage(ctx context.Context, session *types.Session) {
	//TODO implement me
	panic("implement me")
}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{store: storage.NewMemoryStorage("session")}
}
