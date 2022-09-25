package dto

import "gim/internal/module/gateway/domain/types"

type SessionQueryRequest struct {
	SessionType string `json:"session_type"`
}

type SessionQueryResponse struct {
	Items []types.Session `json:"items"`
}

type MessageQueryRequest struct {
	SessionId string `json:"session_id"`
	MessageId string `json:"message_id"`
	Count     int    `json:"count"`
}

type MessageQueryResponse struct {
	Items []types.Message `json:"items"`
}

type MessageSendRequest struct {
	Type        string `json:"type"`
	TargetId    string `json:"target_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	RequestId   string `json:"request_id"`
}
type MessageSendResponse struct{}
