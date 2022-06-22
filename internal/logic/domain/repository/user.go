package repository

import (
	"context"
	"gim/internal/logic/domain/entity"
)

type UserRepo interface {
	Find(ctx context.Context, id string) (*entity.User, error)
	Save(ctx context.Context, item *entity.User) error
}
