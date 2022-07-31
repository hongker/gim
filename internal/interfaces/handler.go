package interfaces

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/pkg/network"
)

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
