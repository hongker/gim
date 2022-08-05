package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/store"
	uuid "github.com/satori/go.uuid"
)

type UserRepo struct {
	store *store.Hash
}

func (repo *UserRepo) Save(ctx context.Context, user *entity.User) (err error) {
	user.Id = uuid.NewV4().String()
	return repo.store.Save(user.Id, user)
}


func NewUserRepo() repository.UserRepository {
	return &UserRepo{store: store.NewHash()}
}


