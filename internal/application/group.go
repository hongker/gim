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
	userRepo repository.UserRepository
}

func (app *GroupApp) Join(ctx context.Context, user *dto.User, groupId string) error {
	group, _ := app.groupRepo.Find(ctx, groupId)
	if group == nil {
		group = &entity.Group{
			Id: groupId,
			Title: fmt.Sprintf("group:%s", groupId),
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

func (app *GroupApp) QueryMember(ctx context.Context, groupId string) (*dto.GroupMemberResponse, error) {
	memberIds, err := app.groupUserRepo.FindAll(ctx, groupId)
	if err != nil {
		return nil, err
	}

	res := &dto.GroupMemberResponse{Items: make([]dto.User, 0, len(memberIds))}
	for _, id := range memberIds {
		user, err := app.userRepo.Find(ctx, id)
		if err != nil {
			continue
		}
		res.Items = append(res.Items, dto.User{
			Id:   user.Id,
			Name: user.Name,
		})
	}

	return res, nil
}
func NewGroupApp(groupRepo repository.GroupRepo, groupUserRepo repository.GroupUserRepo, userRepository repository.UserRepository) (*GroupApp)   {
	return &GroupApp{groupRepo: groupRepo, groupUserRepo: groupUserRepo, userRepo: userRepository}
}
