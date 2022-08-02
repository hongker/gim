package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"sync"
)

type GroupRepo struct {
	mu sync.RWMutex
	items map[string]*entity.Group
}

func (repo *GroupRepo) Find(ctx context.Context, id string) (*entity.Group, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	item, ok := repo.items[id]
	if !ok {
		return nil, errors.DataNotFound("group not found")
	}
	return item, nil
}

func (repo *GroupRepo) Create(ctx context.Context, group *entity.Group) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, ok := repo.items[group.GroupId]
	if ok {
		return errors.DataNotFound("group exist")
	}
	repo.items[group.GroupId] = group
	return nil
}

func NewGroupRepo() repository.GroupRepo {
	return &GroupRepo{
		items: make(map[string]*entity.Group),
	}
}
