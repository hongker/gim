package repository

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/infrastructure/storage"
)

type UserRepository interface {
	Save(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}

type userRepo struct {
	store *storage.StorageManager
}

func (repo *userRepo) Save(ctx context.Context, item *entity.User) error {
	return repo.store.User().Create(ctx, item)
}

func (repo *userRepo) Find(ctx context.Context, id string) (*entity.User, error) {
	return repo.store.User().Find(ctx, id)
}

func NewUserRepository() UserRepository {
	return &userRepo{store: storage.MemoryManager()}
}
