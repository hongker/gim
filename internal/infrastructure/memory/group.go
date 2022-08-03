package memory

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"gim/pkg/store"
)

type GroupRepo struct {
	store *store.Hash
}

func (repo *GroupRepo) Find(ctx context.Context, id string) (*entity.Group, error) {
	item, ok := repo.store.Get(id)
	if !ok {
		return nil, errors.DataNotFound("group not found")
	}
	return item.(*entity.Group), nil
}

func (repo *GroupRepo) Create(ctx context.Context, group *entity.Group) error {
	return repo.store.Save(group.GroupId, group)
}

func NewGroupRepo() repository.GroupRepo {
	return &GroupRepo{
		store: store.NewHash(),
	}
}
