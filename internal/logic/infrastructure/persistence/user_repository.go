package persistence

import (
	"context"
	"gim/internal/logic/domain/entity"
	"gim/internal/logic/domain/repository"
	"github.com/pkg/errors"
	"sync"
)

type UserRepository struct {
	rmu   sync.RWMutex
	users map[string]*entity.User
}

func NewUserRepository() repository.UserRepo {
	return &UserRepository{
		users: make(map[string]*entity.User, 0),
	}
}

func (repo *UserRepository) Find(ctx context.Context, id string) (*entity.User, error) {
	repo.rmu.RLock()
	defer repo.rmu.RUnlock()

	item, ok := repo.users[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return item, nil
}

func (repo *UserRepository) Save(ctx context.Context, item *entity.User) error {
	repo.rmu.Lock()
	repo.users[item.Id] = item
	repo.rmu.Unlock()
	return nil
}
