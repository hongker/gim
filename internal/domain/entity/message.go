package entity

import (
	"gim/api"
	"gim/internal/domain/dto"
)

type Message struct {
	Id          string
	SessionId   string
	SessionType string
	Content     string
	ContentType string
	CreatedAt   int64
	RequestId string
	Sequence    int64
	FromUser    *User
	Group *Group
}

func (item Message) Session() dto.Session  {
	title := item.FromUser.Name
	if item.SessionType == api.GroupSession {
		title = item.Group.Title
	}
	return dto.Session{
		Id:    item.SessionId,
		Type:  item.SessionType,
		Title: title,
	}
}
