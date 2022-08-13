package repository

import (
	"context"
	"gim/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}
