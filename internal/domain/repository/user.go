package repository

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/infra/storage"
	"sync"
)

type UserRepository interface {
	Save(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}

type UserRepo struct {
	mu    sync.RWMutex
	store storage.Storage
}

func (repo *UserRepo) Save(ctx context.Context, item *entity.User) error {
	return repo.store.Save(ctx, item)
}

func (repo *UserRepo) Find(ctx context.Context, id string) (*entity.User, error) {
	user := entity.NewUserWithID(id)
	err := repo.store.Find(ctx, user)
	return user, err
}

var userRepositoryOnce = struct {
	once     sync.Once
	instance UserRepository
}{}

func NewUserRepository() UserRepository {
	userRepositoryOnce.once.Do(func() {
		userRepositoryOnce.instance = &UserRepo{store: storage.NewMemoryStorage("user")}
	})
	return userRepositoryOnce.instance
}
