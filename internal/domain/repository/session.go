package repository

import (
	"context"
	types2 "gim/internal/domain/types"
	"sync"
)

type SessionRepository interface {
	List(ctx context.Context, uid string) ([]types2.Session, error)
	SaveMessage(ctx context.Context, session *types2.Session, msg *types2.Message) error
	QueryMessage(ctx context.Context, session *types2.Session)
}

type sessionRepo struct {
	mu    sync.Mutex
	items map[string][]string
}

func (repo *sessionRepo) List(ctx context.Context, uid string) ([]types2.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *sessionRepo) SaveMessage(ctx context.Context, session *types2.Session, msg *types2.Message) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.items[session.Id]; !ok {
		repo.items[session.Id] = make([]string, 0, 64)
	}
	repo.items[session.Id] = append(repo.items[session.Id], msg.Id)
	return nil
}

func (repo *sessionRepo) QueryMessage(ctx context.Context, session *types2.Session) {
	//TODO implement me
	panic("implement me")
}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{items: map[string][]string{}}
}
