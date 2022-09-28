package storage

import (
	"context"
	"gim/internal/domain/entity"
)

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
	Create(ctx context.Context, item *entity.User) error
	Find(ctx context.Context, id string) (*entity.User, error)
}
type Chatroom interface {
	Create(ctx context.Context) error
	Find(ctx context.Context, id string)
}

type ChatroomMember interface {
	Create(ctx context.Context) error
	Contain(ctx context.Context, id string)
	Remove(ctx context.Context, id string)
}

type Session interface {
	Create(ctx context.Context) error
	List(ctx context.Context)
}

type SessionMessage interface {
	Create(ctx context.Context) error
	List(ctx context.Context)
}
