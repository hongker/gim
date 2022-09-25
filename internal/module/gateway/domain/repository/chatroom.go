package repository

import (
	"context"
	"gim/internal/module/gateway/domain/types"
	"gim/pkg/store"
	"github.com/ebar-go/ego/errors"
	"sync"
)

type ChatroomRepository interface {
	Create(ctx context.Context, chatroom *types.Chatroom) error
	Update(ctx context.Context, chatroom *types.Chatroom) error
	Find(ctx context.Context, id string) (*types.Chatroom, error)
	AddMember(ctx context.Context, chatroom *types.Chatroom, uid string) error
	HasMember(ctx context.Context, chatroom *types.Chatroom, uid string) bool
}

func NewChatroomRepository() ChatroomRepository {
	return &chatroomRepo{
		items:   map[string]*types.Chatroom{},
		members: map[string]store.Set{},
	}
}

type chatroomRepo struct {
	mu      sync.Mutex // guards
	items   map[string]*types.Chatroom
	members map[string]store.Set
}

func (c *chatroomRepo) Create(ctx context.Context, chatroom *types.Chatroom) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[chatroom.Id] = chatroom
	return nil
}

func (c *chatroomRepo) Update(ctx context.Context, chatroom *types.Chatroom) error {
	//TODO implement me
	panic("implement me")
}

func (c *chatroomRepo) Find(ctx context.Context, id string) (*types.Chatroom, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[id]
	if !ok {
		return nil, errors.NotFound("chatroom not found")
	}
	return item, nil
}

func (c *chatroomRepo) AddMember(ctx context.Context, chatroom *types.Chatroom, uid string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.members[chatroom.Id]; !ok {
		c.members[chatroom.Id] = store.ThreadSafe()
	}
	c.members[chatroom.Id].Add(uid)
	return nil
}

func (c *chatroomRepo) HasMember(ctx context.Context, chatroom *types.Chatroom, uid string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.members[chatroom.Id]; !ok {
		return false
	}
	return c.members[chatroom.Id].Contain(uid)
}
