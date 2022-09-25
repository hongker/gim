package repository

import (
	"context"
	"gim/internal/module/gateway/domain/types"
	"sync"
)

type SessionRepository interface {
	List(ctx context.Context, uid string) ([]types.Session, error)
	SaveMessage(ctx context.Context, session *types.Session, msg *types.Message) error
	QueryMessage(ctx context.Context, session *types.Session)
}

type sessionRepo struct {
	mu    sync.Mutex
	items map[string][]string
}

func (repo *sessionRepo) List(ctx context.Context, uid string) ([]types.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *sessionRepo) SaveMessage(ctx context.Context, session *types.Session, msg *types.Message) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.items[session.Id]; !ok {
		repo.items[session.Id] = make([]string, 0, 64)
	}
	repo.items[session.Id] = append(repo.items[session.Id], msg.Id)
	return nil
}

func (repo *sessionRepo) QueryMessage(ctx context.Context, session *types.Session) {
	//TODO implement me
	panic("implement me")
}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{}
}
