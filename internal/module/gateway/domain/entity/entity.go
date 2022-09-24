package entity

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id          string `json:"id"`
	SenderId    string `json:"sender_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Status      int    `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

type SessionMessage struct {
	Uid       string `json:"uid"`
	SessionId string `json:"session_id"`
	MessageId string `json:"message_id"`
}
