package repository

import (
	"context"
	"gim/internal/domain/entity"
)

type GroupRepo interface {
	Find(ctx context.Context, id string) (*entity.Group, error)
	Create(ctx context.Context, group *entity.Group) ( error)
}

type GroupUserRepo interface {
	Create(ctx context.Context, item *entity.GroupUser) error
	Find(ctx context.Context, groupId string, userId string) (*entity.GroupUser, error)
	Delete(ctx context.Context, groupId string, userId string) error
	FindAll(ctx context.Context, groupId string) ([]string, error)
}