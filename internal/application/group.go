package application

import (
	"context"
	"fmt"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"gim/pkg/errors"
	"time"
)

type GroupApp struct {
	groupRepo repository.GroupRepo
	groupUserRepo repository.GroupUserRepo
}

func (app *GroupApp) Join(ctx context.Context, user *dto.User, groupId string) error {
	group, _ := app.groupRepo.Find(ctx, groupId)
	if group == nil {
		group = &entity.Group{
			Id: groupId,
			Title: fmt.Sprintf("group:%d", groupId),
			Creator: user.Id,
			CreatedAt: time.Now().Unix(),
		}
		if err := app.groupRepo.Create(ctx, group); err != nil {
			return errors.WithMessage(err, "create group")
		}
	}

	groupUser, _ := app.groupUserRepo.Find(ctx, group.Id, user.Id)
	if groupUser != nil {
		return errors.Failure("group user is exist")
	}

	if err := app.groupUserRepo.Create(ctx, &entity.GroupUser{
		GroupId:   group.Id,
		UserId:    user.Id,
		CreatedAt: time.Now().Unix(),
	}); err != nil {
		return errors.WithMessage(err, "create group user")
	}

	return nil
}

func (app *GroupApp) Leave(ctx context.Context, user *dto.User, groupId string) error {
	return app.groupUserRepo.Delete(ctx, groupId ,user.Id)
}
func NewGroupApp(groupRepo repository.GroupRepo, groupUserRepo repository.GroupUserRepo) (*GroupApp)   {
	return &GroupApp{groupRepo: groupRepo, groupUserRepo: groupUserRepo}
}
