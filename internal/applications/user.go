package applications

import "gim/internal/domain/dto"

type UserApp struct {

}

func (u *UserApp) Login(req *dto.UserLoginRequest) (res *dto.UserLoginResponse,err error)    {
	return
}
