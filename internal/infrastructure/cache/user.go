package cache

import (
	"context"
	"gim/internal/domain/repository"
)

type UserRepo struct {
	repository.UserRepository
}

func (repo *UserRepo) Save(ctx context.Context) (err error) {
	err = repo.UserRepository.Save(ctx)
	if err != nil {
		return
	}

	return
}

func NewUserRepo(delegate repository.UserRepository) *UserRepo {
	return &UserRepo{UserRepository: delegate}
}


