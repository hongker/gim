package handler

import (
	"gim/internal/applications"
	"gim/internal/domain/dto"
	"gim/internal/interfaces/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
)

type UserHandler struct {
	userApp *applications.UserApp
	gateApp *applications.GateApp
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

	handler.gateApp.RegisterConn(resp.UID, ctx.Connection())

	return resp, nil
}

func NewUserHandler(userApp *applications.UserApp,
gateApp *applications.GateApp) *UserHandler {
	return &UserHandler{
		userApp: userApp,
		gateApp: gateApp,
	}
}
