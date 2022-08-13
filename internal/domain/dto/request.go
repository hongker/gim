package dto

import (
	"fmt"
	"gim/api"
)

type UserLoginRequest struct {
	UID string `json:"uid"`
	Name string `json:"name"`
}

type UserLoginResponse struct {
	UID string `json:"uid"`
	Name string `json:"name"`
}

type MessageSendRequest struct {
	Type        string `json:"type"`
	TargetId string `json:"target_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	RequestId string `json:"request_id"`
}

func (req MessageSendRequest) SessionId(uid string) string {
	if req.Type == api.GroupSession {
		return fmt.Sprintf("%v:%s", req.Type, req.TargetId)
	}

	if uid > req.TargetId {
		return fmt.Sprintf("%v:%s:%s", req.Type, req.TargetId, uid)
	}
	return fmt.Sprintf("%v:%s:%s", req.Type, uid, req.TargetId)
}

type MessageSendResponse struct {}

type MessageQueryRequest struct {
	SessionId string `json:"session_id"`
	Last int64 `json:"last"`
	Limit int `json:"limit"`
}
type MessageQueryResponse struct {
	Items []Message `json:"items"`
}

type Message struct {
	Id string `json:"id"`
	RequestId string `json:"request_id"`
	Session Session `json:"session"`
	ContentType string `json:"content_type"`
	Content string `json:"content"`
	CreatedAt int64 `json:"created_at"`
	Sequence int64 `json:"sequence"`
	FromUser User `json:"from_user"`
}

type GroupJoinRequest struct {
	GroupId string `json:"group_id"`
}
type GroupJoinResponse struct{}

type GroupLeaveRequest struct {
	GroupId string `json:"group_id"`
}
type GroupLeaveResponse struct{}