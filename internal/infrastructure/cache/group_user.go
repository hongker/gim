package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"sync"
)

type GroupUserRepo struct {
	rmu sync.RWMutex
	items map[string]map[string]int64
}

func (repo *GroupUserRepo) Create(ctx context.Context, item *entity.GroupUser) error {
	repo.rmu.Lock()
	if _, ok := repo.items[item.GroupId]; !ok {
		repo.items[item.GroupId] = make(map[string]int64)
	}
	repo.items[item.GroupId][item.UserId] = int64(len(repo.items[item.GroupId])) + 1
	repo.rmu.Unlock()
	return nil
}

func (repo *GroupUserRepo) Find(ctx context.Context, groupId string, userId string) (*entity.GroupUser, error) {
	repo.rmu.RLock()
	defer repo.rmu.RUnlock()
	groupUsers , ok := repo.items[groupId]
	if !ok {
		return nil, errors.DataNotFound("not found")
	}

	createdAt, ok := groupUsers[userId]
	if !ok {
		return nil, errors.DataNotFound("not found")
	}
	return &entity.GroupUser{
		GroupId:   groupId,
		UserId:    userId,
		CreatedAt: createdAt,
	}, nil


}

func (repo *GroupUserRepo) Delete(ctx context.Context, groupId string, userId string) error {
	repo.rmu.Lock()
	defer repo.rmu.Unlock()
	groupUsers , ok := repo.items[groupId]
	if !ok {
		return nil
	}
	delete(groupUsers, userId)

	return nil
}

func NewGroupUserRepo() repository.GroupUserRepo  {
	return &GroupUserRepo{
		items: map[string]map[string]int64{},
	}
}

