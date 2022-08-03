package handler

import (
	"gim/internal/aggregate"
	"gim/internal/domain/dto"
	"gim/internal/domain/event"
	"gim/internal/interfaces/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
)

type UserHandler struct {
	userApp *aggregate.UserApp
}

func (handler *UserHandler) Login(ctx *network.Context) (interface{}, error)  {
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

func NewUserHandler(userApp *aggregate.UserApp,) *UserHandler {
	return &UserHandler{
		userApp: userApp,
	}
}
