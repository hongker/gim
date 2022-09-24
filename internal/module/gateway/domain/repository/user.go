package repository

import (
	"context"
	"gim/internal/module/gateway/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}

type UserRepo struct{}

func (repo UserRepo) Save(ctx context.Context, item *entity.User) error {
	return nil
}

func (repo UserRepo) Find(ctx context.Context, id string) (*entity.User, error) {
	return &entity.User{Id: id}, nil
}

func NewUserRepository() UserRepository {
	return &UserRepo{}
}
