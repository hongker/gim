package aggregate

import (
	"context"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
)

type UserApp struct {
	repo repository.UserRepository
}

func NewUserApp(repo repository.UserRepository) *UserApp {
	return &UserApp{repo: repo}
}

func (app *UserApp) Login(ctx context.Context, req *dto.UserLoginRequest) (res *dto.UserLoginResponse,err error)    {
	user := &entity.User{Name: req.Name}
	if err = app.repo.Save(ctx, user); err != nil {
		return
	}

	res = &dto.UserLoginResponse{UID: user.Id, Name: req.Name}
	return
}
