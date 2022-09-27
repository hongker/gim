package entity

type Primary struct {
	Id string `json:"id"`
}

func (p Primary) ID() string {
	return p.Id
}

type User struct {
	Primary
	Name string `json:"name"`
}

func NewUserWithID(id string) *User {
	return &User{Primary: Primary{Id: id}}
}

type Message struct {
	Primary
	SenderId  string `json:"sender_id"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

func NewMessageWithID(id string) *Message {
	return &Message{Primary: Primary{Id: id}}
}

type SessionMessage struct {
	Uid       string `json:"uid"`
	SessionId string `json:"session_id"`
	MessageId string `json:"message_id"`
}
