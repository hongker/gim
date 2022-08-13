package handler

import (
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/presentation/helper"
	"gim/pkg/errors"
	"gim/pkg/network"
)

type MessageHandler struct {
	messageApp *application.MessageApp
}


func (handler *MessageHandler) Send(ctx *network.Context) (interface{}, error) {
	req := &dto.MessageSendRequest{}
	if err := helper.Bind(ctx, req); err != nil {
		return nil, errors.InvalidParameter(err.Error())
	}

	sender := helper.GetContextUser(ctx)

	err := handler.messageApp.Send(ctx, sender, req)
	if err != nil {
		return nil, errors.WithMessage(err, "send message")
	}


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

func NewMessageHandler(messageApp *application.MessageApp) *MessageHandler {
	return &MessageHandler{
		messageApp: messageApp,
	}
}