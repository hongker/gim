package application

import (
	"context"
	"gim/internal/module/gateway/domain/dto"
	"gim/internal/module/gateway/domain/entity"
	"gim/internal/module/gateway/domain/repository"
	"gim/internal/module/gateway/domain/types"
	"github.com/ebar-go/ego/errors"
)

type UserApplication interface {
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

func NewUserApplication() UserApplication {
	return &userApplication{
		repo: repository.NewUserRepository(),
		auth: types.DefaultAuthenticator(),
	}
}

type userApplication struct {
	repo repository.UserRepository
	auth types.Authenticator
}

func (app userApplication) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user := &entity.User{Id: req.ID, Name: req.Name}
	if err := app.repo.Save(ctx, user); err != nil {
		return nil, errors.WithMessage(err, "save user")
	}

	token, err := app.auth.GenerateToken(ctx, user.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "generate token")
	}

	resp := &dto.UserLoginResponse{Token: token}
	return resp, nil
}

func (app userApplication) Authenticate(ctx context.Context) {

}
