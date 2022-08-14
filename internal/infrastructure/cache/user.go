package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"github.com/patrickmn/go-cache"
)

type UserRepo struct {
	repository.UserRepository
	store *cache.Cache
}

func (repo *UserRepo) Find(ctx context.Context, id string) (*entity.User, error)   {
	item, ok := repo.store.Get(id)
	if ok {
		return item.(*entity.User), nil
	}
	if repo.UserRepository == nil {
		return nil, errors.DataNotFound("user not found")
	}
	res, err := repo.UserRepository.Find(ctx, id)
	if err != nil {
		return nil, errors.DataNotFound("user not found")
	}

	repo.store.Set(id, res, cache.DefaultExpiration)
	return res, nil
}

func (repo *UserRepo) Save(ctx context.Context, user *entity.User) ( error) {
	if repo.UserRepository != nil {
		return repo.UserRepository.Save(ctx, user)
	}

	repo.store.Set(user.Id, user, cache.NoExpiration)
	return nil
}


func NewUserRepo(delegate repository.UserRepository,store *cache.Cache) repository.UserRepository {
	return &UserRepo{UserRepository: delegate,store: store}
}



