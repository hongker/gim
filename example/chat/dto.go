package main

import "encoding/json"

type LoginRequest struct {
	Name string `json:"name"`
}
type LoginResponse struct {
	ID string `json:"id"`
}

type SubscribeChannelRequest struct {
	ID string `json:"id"`
}
type SubscribeChannelResponse struct{}

type SendMessageRequest struct {
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

type SendMessageResponse struct {
	MsgID string `json:"msg_id"`
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
