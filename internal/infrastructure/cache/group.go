package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/internal/infrastructure/config"
	"gim/pkg/errors"
	"github.com/patrickmn/go-cache"
)

type GroupRepo struct {
	repository.GroupRepo
	store *cache.Cache
}

func (repo *GroupRepo) Find(ctx context.Context, id string) (*entity.Group, error) {
	item, ok := repo.store.Get(id)
	if ok {
		return item.(*entity.Group), nil
	}
	res, err := repo.GroupRepo.Find(ctx, id)
	if err != nil {
		return nil, errors.DataNotFound("group not found")
	}

	repo.store.Set(id, res, cache.DefaultExpiration)
	return res, nil

}


func NewGroupRepo(delegate repository.GroupRepo, conf *config.Config) repository.GroupRepo {
	return &GroupRepo{
		GroupRepo: delegate,
		store: cache.New(conf.Cache.Expired, conf.Cache.Cleanup),
	}
}
