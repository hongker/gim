package types

import (
	"encoding/json"
	"gim/internal/domain/entity"
	uuid "github.com/satori/go.uuid"
	"time"
)

type MessageCategory string

const (
	MessageText    MessageCategory = "text"
	MessagePicture MessageCategory = "picture"
)

const (
	MessageStatusNormal = iota + 1
)

type Message struct {
	Id        string          `json:"id"`
	SenderId  string          `json:"sender_id"`
	Category  MessageCategory `json:"category"`
	Content   string          `json:"content"`
	Sequence  int64           `json:"sequence"`
	Status    int             `json:"status"`
	CreatedAt int64           `json:"created_at"`
}

func (msg *Message) Entity() *entity.Message {
	return &entity.Message{
		Id:        msg.Id,
		SenderId:  msg.SenderId,
		Content:   msg.Content,
		Category:  string(msg.Category),
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
	}
}

func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

func NewMessage(senderId string, category MessageCategory, content string) *Message {
	return &Message{
		Id:        uuid.NewV4().String(),
		SenderId:  senderId,
		Category:  category,
		Content:   content,
		Status:    MessageStatusNormal,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func NewTextMessage(senderId string, content string) *Message {
	return NewMessage(senderId, MessageText, content)
}

func NewPictureMessage(senderId string, content string) *Message {
	return NewMessage(senderId, MessagePicture, content)
}

type SessionMessage struct {
	Session *Session
	Message *Message
}

type SessionMessageItems struct {
	Session *Session   `json:"session"`
	Items   []*Message `json:"items"`
}
