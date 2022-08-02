package dto

type UserLoginRequest struct {
	Name string `json:"name"`
}

type UserLoginResponse struct {
	UID string `json:"uid"`
	Name string `json:"name"`
}

type MessageSendRequest struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	ClientMsgId string `json:"client_msg_id"`
	SessionId   string `json:"session_id"`

}

type MessageSendResponse struct {}

type MessageQueryRequest struct {
	SessionId string `json:"session_id"`
	Last int64 `json:"last"`
}
type MessageQueryResponse struct {
	Items []Message `json:"items"`
}

type Message struct {
	SessionId string `json:"session_id"`
	Content string `json:"content"`
	CreatedAt int64 `json:"created_at"`
}

type GroupJoinRequest struct {
	GroupId string `json:"group_id"`
}
type GroupJoinResponse struct{}

type GroupLeaveRequest struct {
	GroupId string `json:"group_id"`
}
type GroupLeaveResponse struct{}