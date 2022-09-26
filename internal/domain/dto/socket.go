package dto

type SocketHeartbeatRequest struct{}
type SocketHeartbeatResponse struct {
	ServerTime int64 `json:"server_time"`
}
