package main

import "encoding/json"

type LoginRequest struct{ Name string }
type LoginResponse struct {
	ID string
}

type SubscribeChannelRequest struct{ ID string }
type SubscribeChannelResponse struct{}

type SendMessageRequest struct {
	ChannelID string
	Content   string
}

type SendMessageResponse struct {
	MsgID string
}

type Message struct {
	ID      string      `json:"id"`
	Content string      `json:"content"`
	Sender  MessageUser `json:"sender"`
}

func (message Message) Serialize() []byte {
	bytes, _ := json.Marshal(message)
	return bytes
}

type MessageUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
