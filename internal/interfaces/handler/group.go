package handler

import (
	"gim/internal/aggregate"
	"gim/internal/domain/dto"
	"gim/internal/domain/event"
	"gim/internal/interfaces/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
)

type GroupHandler struct {
	groupApp *aggregate.GroupApp
}

func (handler *GroupHandler) Join(ctx *network.Context) (interface{}, error) {
	req := &dto.GroupJoinRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	user := helper.GetContextUser(ctx)
	err := handler.groupApp.Join(ctx, user, req.GroupId)
	if err != nil {
		return nil, errors.WithMessage(err, "join group")
	}

	event.Trigger(event.JoinGroup, req.GroupId, ctx.Connection())

	return nil, nil
}

func (handler *GroupHandler) Leave(ctx *network.Context) ( interface{},  error)  {
	req := &dto.GroupLeaveRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	user := helper.GetContextUser(ctx)
	err := handler.groupApp.Leave(ctx, user, req.GroupId)
	if err != nil {
		return nil, errors.WithMessage(err, "leave group")
	}
	event.Trigger(event.LeaveGroup, req.GroupId, ctx.Connection())
	return nil, nil
}

func NewGroupHandler(groupApp *aggregate.GroupApp, ) *GroupHandler {
	return &GroupHandler{groupApp: groupApp, }
}