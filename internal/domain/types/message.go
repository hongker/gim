package types

import (
	"encoding/json"
	"gim/internal/domain/entity"
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

func NewMessage(category MessageCategory, content string) *Message {
	return &Message{
		Category:  category,
		Content:   content,
		Status:    MessageStatusNormal,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func NewTextMessage(content string) *Message {
	return NewMessage(MessageText, content)
}

func NewPictureMessage(content string) *Message {
	return NewMessage(MessagePicture, content)
}

type SessionMessage struct {
	Session *Session
	Message *Message
}

type MessagePacket struct {
	Session *Session   `json:"session"`
	Items   []*Message `json:"items"`
}

func (packet MessagePacket) Encode() []byte {
	b, _ := json.Marshal(packet)
	return b
}
