package types

import (
	"fmt"
	"gim/internal/domain/entity"
	"strconv"
	"strings"
)

var (
	SessionChatroom = 1
	SessionPrivate  = 2
)

func SessionId(category int, targetId string) string {
	return fmt.Sprintf("%v:%s", category, targetId)
}

func CategoryFromSessionId(id string) int {
	category, _ := strconv.Atoi(strings.Split(id, ":")[0])
	return category
}

type Session struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (s Session) Entity() *entity.Session {
	return &entity.Session{
		Id:    s.Id,
		Title: s.Title,
	}
}

func (s Session) Category() int {
	return CategoryFromSessionId(s.Id)
}

func (s Session) IsChatroom() bool {
	return s.Category() == SessionChatroom
}
func (s Session) IsPrivate() bool {
	return s.Category() == SessionPrivate
}

func (s Session) GetPrivateUid() string {
	return strings.Split(s.Id, ":")[1]
}
func (s Session) GetChatroomId() string {
	return strings.Split(s.Id, ":")[1]
}

func NewSession(id string, title string) *Session {
	return &Session{id, title}
}

func NewChatroomSession(roomId string, title string) *Session {
	return NewSession(SessionId(SessionChatroom, roomId), title)
}

func NewPrivateSession(senderId, receiverId string, title string) *Session {
	userIds := []string{senderId, receiverId}
	targetId := strings.Join(userIds, ":")
	return NewSession(SessionId(SessionPrivate, targetId), title)
}
