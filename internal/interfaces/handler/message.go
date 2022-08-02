package handler

import (
	"gim/api"
	"gim/internal/applications"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/pkg/network"
)

type MessageHandler struct {
	messageApp *applications.MessageApp
	gateApp *applications.GateApp
}


func (handler *MessageHandler) Send(ctx *network.Context, p *api.Packet) error {
	req := &dto.MessageSendRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	fromUser := &entity.User{}

	message, err := handler.messageApp.Send(ctx, fromUser, req)
	if err != nil {
		return err
	}

	packet := api.BuildPacket(api.OperateMessagePush, message)
	handler.gateApp.Push(req.Type, req.SessionId, packet.Encode())

	return p.Marshal(nil)
}


func (handler *MessageHandler) Query(ctx *network.Context, p *api.Packet) error {
	req := &dto.MessageQueryRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	resp, err := handler.messageApp.Query(ctx, req)
	if err != nil {
		return err
	}
	return p.Marshal(resp)
}

func NewMessageHandler(messageApp *applications.MessageApp,
gateApp *applications.GateApp) *MessageHandler {
	return &MessageHandler{
		messageApp: messageApp,
		gateApp:    gateApp,
	}
}