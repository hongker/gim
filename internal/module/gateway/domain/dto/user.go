package dto

import "github.com/ebar-go/ego/errors"

type UserLoginRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (req UserLoginRequest) Validate() error {
	if req.ID == "" {
		return errors.InvalidParam("invalid id")
	}

	return nil
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserFindRequest struct{}
type UserFindResponse struct{}
