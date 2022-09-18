package application

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
)

type UserApplication interface {
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

func NewUserApplication() UserApplication {
	return &userApplication{}
}

type userApplication struct{}

func (userApplication) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	return nil, nil
}
