package handler

import (
	"gim/api"
	"gim/internal/aggregate"
	"gim/internal/domain/dto"
	"gim/internal/interfaces/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
)

type MessageHandler struct {
	messageApp *aggregate.MessageApp
	gateApp *aggregate.GateApp
}


func (handler *MessageHandler) Send(ctx *network.Context) (interface{}, error) {
	req := &dto.MessageSendRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	sender := helper.GetContextUser(ctx)

	message, err := handler.messageApp.Send(ctx, sender, req)
	if err != nil {
		return nil, errors.WithMessage(err, "send message")
	}

	packet := api.BuildPacket(api.OperateMessagePush, message)
	handler.gateApp.Push(req.Type, req.SessionId, packet.Encode())

	return nil, nil
}


func (handler *MessageHandler) Query(ctx *network.Context) (interface{}, error) {
	req := &dto.MessageQueryRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	resp, err := handler.messageApp.Query(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "query message")
	}
	return resp, nil
}

func NewMessageHandler(messageApp *aggregate.MessageApp,
gateApp *aggregate.GateApp) *MessageHandler {
	return &MessageHandler{
		messageApp: messageApp,
		gateApp:    gateApp,
	}
}