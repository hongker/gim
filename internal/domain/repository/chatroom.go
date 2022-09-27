package repository

import (
	"context"
	"gim/internal/domain/types"
	"gim/internal/infra/storage"
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
			store:       storage.NewMemoryStorage("chatroom"),
			memberStore: storage.NewMemoryStorage("chatroomMember"),
		}
	})
	return chatroomRepoOnce.instance
}

type chatroomRepo struct {
	mu          sync.Mutex // guards
	store       storage.Storage
	memberStore storage.Storage
}

func (repo *chatroomRepo) GetMember(ctx context.Context, roomId string) ([]string, error) {
	return nil, nil
}

func (repo *chatroomRepo) Create(ctx context.Context, chatroom *types.Chatroom) error {
	return repo.store.Save(ctx, chatroom)
}

func (repo *chatroomRepo) Update(ctx context.Context, chatroom *types.Chatroom) error {
	return repo.store.Save(ctx, chatroom)
}

func (repo *chatroomRepo) Find(ctx context.Context, id string) (*types.Chatroom, error) {
	item := &types.Chatroom{Id: id}
	if err := repo.store.Find(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *chatroomRepo) AddMember(ctx context.Context, chatroom *types.Chatroom, uid string) error {
	cm := &types.ChatroomMember{Id: chatroom.Id}
	if err := repo.memberStore.Find(ctx, cm); err != nil {
		if !errors.Is(err, errors.NotFound("")) {
			return err
		}

	}
	if cm.Members == nil {
		cm.Members = store.ThreadSafe()
	}
	cm.Members.Add(uid)
	return repo.memberStore.Save(ctx, cm)
}

func (repo *chatroomRepo) HasMember(ctx context.Context, chatroom *types.Chatroom, uid string) bool {
	cm := &types.ChatroomMember{Id: chatroom.Id}
	if err := repo.memberStore.Find(ctx, cm); err != nil {
		return false
	}
	if cm.Members == nil {
		return false
	}
	return cm.Members.Contain(uid)
}
