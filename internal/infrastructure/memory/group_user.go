package memory

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"gim/pkg/store"
	"sync"
)

type GroupUserRepo struct {
	rmu sync.RWMutex
	items map[string]store.Set
}

func (repo *GroupUserRepo) Create(ctx context.Context, item *entity.GroupUser) error {
	repo.rmu.Lock()
	if _, ok := repo.items[item.GroupId]; !ok {
		repo.items[item.GroupId] = store.ThreadSafe()
	}
	repo.items[item.GroupId].Add(item.UserId)
	repo.rmu.Unlock()
	return nil
}

func (repo *GroupUserRepo) Find(ctx context.Context, groupId string, userId string) (*entity.GroupUser, error) {
	repo.rmu.RLock()
	defer repo.rmu.RUnlock()
	set , ok := repo.items[groupId]
	if !ok || !set.Contain(userId){
		return nil, errors.DataNotFound("not found")
	}

	return &entity.GroupUser{
		GroupId:   groupId,
		UserId:    userId,
	}, nil


}

func (repo *GroupUserRepo) Delete(ctx context.Context, groupId string, userId string) error {
	repo.rmu.Lock()
	defer repo.rmu.Unlock()
	set , ok := repo.items[groupId]
	if !ok {
		return nil
	}
	set.Remove(userId)

	return nil
}

func NewGroupUserRepo() repository.GroupUserRepo  {
	return &GroupUserRepo{
		items: map[string]store.Set{},
	}
}

