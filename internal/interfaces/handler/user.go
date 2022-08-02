package handler

import (
	"gim/api"
	"gim/internal/applications"
	"gim/internal/domain/dto"
	"gim/pkg/network"
)

type UserHandler struct {
	userApp *applications.UserApp
	gateApp *applications.GateApp
}

func (handler *UserHandler) Login(ctx *network.Context, p *api.Packet) error  {
	req := &dto.UserLoginRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	resp, err := handler.userApp.Login(ctx, req)
	if err != nil {
		return err
	}

	handler.gateApp.RegisterConn(resp.UID, ctx.Connection())

	return p.Marshal(resp)
}

func NewUserHandler(userApp *applications.UserApp,
gateApp *applications.GateApp) *UserHandler {
	return &UserHandler{
		userApp: userApp,
		gateApp: gateApp,
	}
}
