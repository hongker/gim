package repository

import (
	"context"
	"gim/internal/module/gateway/domain/entity"
	"github.com/ebar-go/ego/errors"
	"sync"
)

type UserRepository interface {
	Save(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}

type UserRepo struct {
	mu    sync.RWMutex
	items map[string]string
}

func (repo *UserRepo) Save(ctx context.Context, item *entity.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.items[item.Id] = item.Name
	return nil
}

func (repo *UserRepo) Find(ctx context.Context, id string) (*entity.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	name, ok := repo.items[id]
	if !ok {
		return nil, errors.NotFound("user not found")
	}

	return &entity.User{Id: id, Name: name}, nil
}

func NewUserRepository() UserRepository {
	return &UserRepo{}
}
