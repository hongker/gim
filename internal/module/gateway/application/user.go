package application

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/errors"
)

type UserApplication interface {
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

func NewUserApplication() UserApplication {
	return &userApplication{}
}

type userApplication struct{}

func (userApplication) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	if req.Name == "" {
		return nil, errors.InvalidParam("invalid name")
	}

	return nil, nil
}
