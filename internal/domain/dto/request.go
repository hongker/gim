package dto

type UserLoginRequest struct {
	Name string `json:"name"`
}

type UserLoginResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
}