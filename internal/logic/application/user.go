package application

import (
	"context"
	"gim/internal/logic/domain/entity"
	"gim/internal/logic/domain/repository"
)

type UserApp struct {
	userRepo repository.UserRepo
}

func NewUserApp(userRepo repository.UserRepo) *UserApp {
	return &UserApp{userRepo: userRepo}
}

func (app *UserApp) Find(id string) (user *entity.User, err error) {
	return
}

func (app *UserApp) Auth(ctx context.Context, id string, name string) (err error) {
	err = app.userRepo.Save(ctx, &entity.User{
		Id:   id,
		Name: name,
	})
	return
}
