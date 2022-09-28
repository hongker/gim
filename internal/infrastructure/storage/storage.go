package storage

import "context"

type Storage interface {
	Message() Message
	User() User
	Chatroom() Chatroom
	ChatroomMember() ChatroomMember
}

type Message interface {
	Create(ctx context.Context) error
	Find(ctx context.Context, id string)
}
type User interface {
	Create(ctx context.Context) error
	Find(ctx context.Context, id string)
}
type Chatroom interface {
	Create(ctx context.Context) error
	Find(ctx context.Context, id string)
}

type ChatroomMember interface {
	Create(ctx context.Context) error
	Has(ctx context.Context, id string)
	Remove(ctx context.Context, id string)
}
