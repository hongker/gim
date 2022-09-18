package types

import (
	"fmt"
	"strings"
)

type SessionCategory string

var (
	SessionChatroom SessionCategory = "chatroom"
	SessionPrivate  SessionCategory = "private"
)

func SessionId(category SessionCategory, targetId string) string {
	return fmt.Sprintf("%s:%s", category, targetId)
}

type Session struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (s Session) Category() SessionCategory {
	return SessionCategory(strings.Split(s.Id, ":")[0])
}

func (s Session) IsChatroom() bool {
	return s.Category() == SessionChatroom
}
func (s Session) IsPrivate() bool {
	return s.Category() == SessionPrivate
}

func NewSession(id string, title string) *Session {
	return &Session{id, title}
}

func NewChatroomSession(roomId string, title string) *Session {
	return NewSession(SessionId(SessionChatroom, roomId), title)
}

func NewPrivateSession(userId string, title string) *Session {
	return NewSession(SessionId(SessionPrivate, userId), title)
}
