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
	GetMember(ctx context.Context, roomId string) ([]string, error)
}

var chatroomRepoOnce = struct {
	once     sync.Once
	instance ChatroomRepository
}{}

func NewChatroomRepository() ChatroomRepository {
	chatroomRepoOnce.once.Do(func() {
		chatroomRepoOnce.instance = &chatroomRepo{
			items:   map[string]*types.Chatroom{},
			members: map[string]store.Set{},
		}
	})
	return chatroomRepoOnce.instance
}

type chatroomRepo struct {
	mu      sync.Mutex // guards
	items   map[string]*types.Chatroom
	members map[string]store.Set
}

func (c *chatroomRepo) GetMember(ctx context.Context, roomId string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := make([]string, 0)
	if _, ok := c.members[roomId]; !ok {
		return res, nil
	}
	items := c.members[roomId].ToSlice()
	for _, item := range items {
		res = append(res, item.(string))
	}
	return res, nil
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
