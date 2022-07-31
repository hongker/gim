package applications

import (
	"gim/internal/domain/dto"
	uuid "github.com/satori/go.uuid"
)

type UserApp struct {

}

func NewUserApp() *UserApp {
	return &UserApp{}
}

func (u *UserApp) Login(req *dto.UserLoginRequest) (res *dto.UserLoginResponse,err error)    {
	res = &dto.UserLoginResponse{Id: uuid.NewV4().String(), Name: req.Name}
	return
}
