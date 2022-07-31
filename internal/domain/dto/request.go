package dto

type UserLoginRequest struct {
	Name string `json:"name"`
}

type UserLoginResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type MessageSendRequest struct {}

type MessageSendResponse struct {}

type MessageQueryRequest struct {}
type MessageQueryResponse struct {}