package interfaces

import (
	"gim/api"
	"gim/internal/domain/dto"
)

func (s *Socket) login(p *api.Packet) error  {
	req := &dto.UserLoginRequest{}
	if err := p.Bind(req); err != nil {
		return err
	}

	resp, err := s.userApp.Login(req)
	if err != nil {
		return err
	}


	return p.Marshal(resp)
}
