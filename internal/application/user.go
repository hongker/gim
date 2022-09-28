package application

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/internal/domain/types"
	"github.com/ebar-go/ego/errors"
)

type UserApplication interface {
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	Logout(ctx context.Context, req *dto.UserLogoutRequest) (*dto.UserLogoutResponse, error)
	Find(ctx context.Context, req *dto.UserFindRequest) (*dto.UserFindResponse, error)
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

func (app userApplication) Logout(ctx context.Context, req *dto.UserLogoutRequest) (*dto.UserLogoutResponse, error) {
	return nil, nil
}

func (app userApplication) Authenticate(ctx context.Context) {

}

func (app userApplication) Find(ctx context.Context, req *dto.UserFindRequest) (*dto.UserFindResponse, error) {
	user, err := app.repo.Find(ctx, req.ID)
	if err != nil {
		return nil, errors.WithMessage(err, "find user")
	}
	resp := &dto.UserFindResponse{Name: user.Name}
	return resp, nil
}
