package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	uuid "github.com/satori/go.uuid"
)

type UserRepo struct {}

func (repo *UserRepo) Save(ctx context.Context, user *entity.User) (err error) {
	user.Id = uuid.NewV4().String()
	return
}

func NewUserRepo() repository.UserRepository {
	return &UserRepo{}
}


