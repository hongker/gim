package types

import "gim/pkg/store"

type Chatroom struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	CreatedAt int64  `json:"created_at"`
}

func (c Chatroom) ID() string {
	return c.Id
}

type ChatroomMember struct {
	Id      string    `json:"id"`
	Members store.Set `json:"members"`
}

func (c ChatroomMember) ID() string {
	return c.Id
}

func (c ChatroomMember) AddMember(member string) {
	if c.Members == nil {
		c.Members = store.ThreadSafe()
	}
	c.Members.Add(member)
}
