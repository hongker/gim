package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
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

	if repo.GroupRepo == nil {
		return nil, errors.DataNotFound("group not found")
	}
	res, err := repo.GroupRepo.Find(ctx, id)
	if err != nil {
		return nil, errors.DataNotFound("group not found")
	}

	repo.store.Set(id, res, cache.DefaultExpiration)
	return res, nil

}

func (repo *GroupRepo) Create(ctx context.Context, item *entity.Group) (error)  {
	if repo.GroupRepo != nil {
		return repo.GroupRepo.Create(ctx, item)
	}

	repo.store.Set(item.Id, item, cache.DefaultExpiration)
	return nil
}


func NewGroupRepo(delegate repository.GroupRepo, store *cache.Cache) repository.GroupRepo {
	return &GroupRepo{
		GroupRepo: delegate,
		store: store,
	}
}
