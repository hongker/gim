package aggregate

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
			GroupId: groupId,
			Title: fmt.Sprintf("group:%d", groupId),
			CreatedAt: time.Now().Unix(),
		}
		if err := app.groupRepo.Create(ctx, group); err != nil {
			return errors.WithMessage(err, "create group")
		}
	}

	groupUser, _ := app.groupUserRepo.Find(ctx, group.GroupId, user.Id)
	if groupUser != nil {
		return errors.Failure("find group user")
	}

	if err := app.groupUserRepo.Create(ctx, &entity.GroupUser{
		GroupId:   group.GroupId,
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
