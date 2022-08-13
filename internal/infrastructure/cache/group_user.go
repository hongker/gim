package cache

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"gim/pkg/store"
	"sync"
)

type GroupUserRepo struct {
	repository.GroupUserRepo
	rmu sync.RWMutex
	items map[string]store.Set
}

func (repo *GroupUserRepo) Create(ctx context.Context, item *entity.GroupUser) error {
	if repo.GroupUserRepo != nil {
		if err := repo.GroupUserRepo.Create(ctx, item); err != nil {
			return err
		}
	}
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
	if !ok {
		if repo.GroupUserRepo == nil {
			return nil, errors.DataNotFound("not found")
		}
		users, err := repo.GroupUserRepo.FindAll(ctx, groupId)
		if err != nil {
			return nil, err
		}
		set = store.ThreadSafe()
		for _, user := range users {
			set.Add(user)
		}
		repo.items[groupId] = set
	}

	if !set.Contain(userId) {
		return nil, errors.DataNotFound("not found")
	}

	return &entity.GroupUser{
		GroupId:   groupId,
		UserId:    userId,
	}, nil


}

func (repo *GroupUserRepo) Delete(ctx context.Context, groupId string, userId string) error {
	if repo.GroupUserRepo != nil {
		if err := repo.GroupUserRepo.Delete(ctx, groupId, userId); err != nil {
			return err
		}
	}
	repo.rmu.Lock()
	defer repo.rmu.Unlock()
	set , ok := repo.items[groupId]
	if !ok {
		return nil
	}
	set.Remove(userId)

	return nil
}

func (repo *GroupUserRepo) FindAll(ctx context.Context, groupId string) ([]string, error) {
	if repo.GroupUserRepo != nil {
		return repo.GroupUserRepo.FindAll(ctx, groupId)
	}

	repo.rmu.RLock()
	defer repo.rmu.RUnlock()
	set , ok := repo.items[groupId]
	if !ok {
		return nil, errors.DataNotFound("group is empty")
	}

	items := set.ToSlice()
	res := make([]string,  len(items))
	for i, item := range items {
		res[i] = item.(string)
	}
	return res, nil
}


func NewGroupUserRepo(delegate repository.GroupUserRepo) repository.GroupUserRepo  {
	return &GroupUserRepo{
		GroupUserRepo: delegate,
		items: map[string]store.Set{},
	}
}

