package dto

type SocketLoginRequest struct {
	Token string `json:"token"`
}
type SocketLoginResponse struct{}

type SocketHeartbeatRequest struct{}
type SocketHeartbeatResponse struct {
	ServerTime int64 `json:"server_time"`
}
