package dto

type SessionQueryRequest struct {
	SessionType string `json:"session_type"`
}

type SessionQueryResponse struct {
	Items []Session `json:"items"`
}

type Session struct {
	Id    string       `json:"id"`
	Title string       `json:"title"`
	Type  int          `json:"type"`
	Last  *MessageItem `json:"last"`
}

type MessageQueryRequest struct {
	SessionId string `json:"session_id"`
	MessageId string `json:"message_id"`
	Count     int    `json:"count"`
}

type MessageQueryResponse struct {
	Items []MessageItem `json:"items"`
}
type MessageUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type MessageItem struct {
	Id      string      `json:"id"`
	Content string      `json:"content"`
	Sender  MessageUser `json:"sender"`
}

type MessageSendRequest struct {
	Type      int    `json:"type"`
	TargetId  string `json:"target_id"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	RequestId string `json:"request_id"`
}
type MessageSendResponse struct{}
