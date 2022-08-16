package handler

import (
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/event"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
	"time"
)

type UserHandler struct {
	userApp *application.UserApp
}

func (handler *UserHandler) Login(ctx *network.Context) (interface{}, error) {
	req := &dto.UserLoginRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	resp, err := handler.userApp.Login(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "login")
	}

	event.Trigger(event.Login, resp.UID, ctx.Connection())

	return resp, nil
}

func (handler *UserHandler) Heartbeat(ctx *network.Context) (interface{}, error) {
	req := &dto.UserHeartbeatRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	event.Trigger(event.Heartbeat, ctx.Connection())
	resp := &dto.UserHeartbeatResponse{Time: time.Now().UnixNano()}
	return resp, nil
}

func NewUserHandler(userApp *application.UserApp) *UserHandler {
	return &UserHandler{
		userApp: userApp,
	}
}
