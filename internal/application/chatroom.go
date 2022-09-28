package application

import (
	"context"
	"fmt"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/internal/domain/repository"
	"github.com/ebar-go/ego/errors"
	"sync"
	"time"
)

type ChatroomApplication interface {
	Create(ctx context.Context, uid string, req *dto.ChatroomCreateRequest) (resp *dto.ChatroomCreateResponse, err error)
	Update(ctx context.Context, uid string, req *dto.ChatroomUpdateRequest) (resp *dto.ChatroomUpdateResponse, err error)
	Join(ctx context.Context, uid string, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error)
	Leave(ctx context.Context, uid string, req *dto.ChatroomLeaveRequest) (resp *dto.ChatroomLeaveResponse, err error)
}

func NewChatroomApplication() ChatroomApplication {
	return &chatroomApplication{
		repo:     repository.NewChatroomRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

type chatroomApplication struct {
	mu       sync.Mutex // guards
	repo     repository.ChatroomRepository
	userRepo repository.UserRepository
}

func (app *chatroomApplication) Create(ctx context.Context, uid string, req *dto.ChatroomCreateRequest) (resp *dto.ChatroomCreateResponse, err error) {
	chatroom := &entity.Chatroom{
		Id:        req.Id,
		Name:      req.Name,
		Creator:   uid,
		CreatedAt: time.Now().UnixMilli(),
	}

	err = app.repo.Create(ctx, chatroom)
	return
}

func (app *chatroomApplication) Update(ctx context.Context, uid string, req *dto.ChatroomUpdateRequest) (resp *dto.ChatroomUpdateResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (app *chatroomApplication) Join(ctx context.Context, uid string, req *dto.ChatroomJoinRequest) (resp *dto.ChatroomJoinResponse, err error) {
	user, err := app.userRepo.Find(ctx, uid)
	if err != nil {
		return nil, errors.WithMessage(err, "find user")
	}
	chatroom, err := app.repo.Find(ctx, req.Id)
	if err != nil && !errors.Is(err, errors.NotFound("")) {
		err = errors.WithMessage(err, "find chatroom")
		return
	}

	if chatroom == nil { // if not exist, create it
		chatroom = &entity.Chatroom{
			Id:        req.Id,
			Name:      fmt.Sprintf("chatroom:%s", req.Id),
			Creator:   uid,
			CreatedAt: time.Now().UnixMilli(),
		}
		if err = app.repo.Create(ctx, chatroom); err != nil {
			return
		}
	}

	err = app.repo.AddMember(ctx, chatroom, user)
	if err != nil {
		return nil, errors.WithMessage(err, "add member")
	}
	return &dto.ChatroomJoinResponse{}, nil
}

func (app *chatroomApplication) Leave(ctx context.Context, uid string, req *dto.ChatroomLeaveRequest) (resp *dto.ChatroomLeaveResponse, err error) {
	//TODO implement me
	panic("implement me")
}
