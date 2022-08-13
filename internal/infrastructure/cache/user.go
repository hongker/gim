package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/internal/infrastructure/config"
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
	res, err := repo.UserRepository.Find(ctx, id)
	if err != nil {
		return nil, errors.DataNotFound("user not found")
	}

	repo.store.Set(id, res, cache.DefaultExpiration)
	return res, nil
}


func NewUserRepo(delegate repository.UserRepository,conf *config.Config) repository.UserRepository {
	return &UserRepo{UserRepository: delegate,store: cache.New(conf.Cache.Expired, conf.Cache.Cleanup)}
}


