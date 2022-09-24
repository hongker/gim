package types

import (
	"fmt"
	"gim/internal/domain/dto"
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

func (s Session) Transform() dto.Session {
	return dto.Session{
		Id:    s.Id,
		Title: s.Title,
		Type:  string(s.Category()),
	}
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

func NewPrivateSession(senderId, receiverId string, title string) *Session {
	userIds := []string{senderId, receiverId}
	//if receiverId > senderId {
	//	userIds[1], userIds[0] = senderId, receiverId
	//}
	targetId := strings.Join(userIds, ":")
	return NewSession(SessionId(SessionPrivate, targetId), title)
}
