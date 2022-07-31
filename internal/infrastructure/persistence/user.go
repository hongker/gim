package persistence

import (
	"context"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func (u *UserRepo) Save(ctx context.Context) error {
	return nil
}

func NewUserRepo(db *gorm.DB) *UserRepo   {
	return &UserRepo{db: db}
}
