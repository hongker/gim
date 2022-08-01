package interfaces

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/entity"
	"gim/pkg/errors"
	"gim/pkg/network"
)

func (s *Socket) WrapHandler(fn func(ctx *network.Context, p *api.Packet) error) Handler{
	return func(ctx *network.Context, p *api.Packet) {
		if err := fn(ctx, p); err != nil {
			Failure(ctx, errors.Convert(err))
		}else {
			Success(ctx, p.Encode())
		}
	}
}
func (s *Socket) login(ctx *network.Context, p *api.Packet) error  {
	req := &dto.UserLoginRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	resp, err := s.userApp.Login(req)
	if err != nil {
		return err
	}

	s.gateApp.RegisterConn(resp.Id, ctx.Connection())

	return p.Marshal(resp)
}

func (s *Socket) send(ctx *network.Context, p *api.Packet) error {
	req := &dto.MessageSendRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	fromUser := &entity.User{}

	message, err := s.messageApp.Send(ctx, fromUser, req)
	if err != nil {
		return err
	}

	packet := api.NewPacket()
	packet.Op = api.OperateMessagePush
	packet.Marshal(message)
	if req.Type == api.PrivateMessage {
		s.gateApp.PushUser(req.SessionId, packet.Encode())
	}else {
		s.gateApp.PushRoom(req.SessionId, packet.Encode())
	}


	return p.Marshal(nil)
}


func (s *Socket) query(ctx *network.Context, p *api.Packet) error {
	req := &dto.MessageQueryRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	resp, err := s.messageApp.Query(ctx, req)
	if err != nil {
		return err
	}
	return p.Marshal(resp)
}

