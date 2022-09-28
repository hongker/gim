package repository

import (
	"context"
	"gim/internal/domain/entity"
	"gim/internal/infrastructure/storage"
)

type ChatroomRepository interface {
	Create(ctx context.Context, chatroom *entity.Chatroom) error
	Find(ctx context.Context, id string) (*entity.Chatroom, error)
	AddMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) error
	HasMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) bool
	RemoveMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) error
}

func NewChatroomRepository() ChatroomRepository {
	return &chatroomRepo{
		store: storage.MemoryManager(),
	}
}

type chatroomRepo struct {
	store *storage.StorageManager
}

func (repo *chatroomRepo) RemoveMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) error {
	return repo.store.Chatroom().RemoveMember(ctx, chatroom.Id, member)
}

func (repo *chatroomRepo) Create(ctx context.Context, chatroom *entity.Chatroom) error {
	return repo.store.Chatroom().Create(ctx, chatroom)
}

func (repo *chatroomRepo) Find(ctx context.Context, id string) (*entity.Chatroom, error) {
	return repo.store.Chatroom().Find(ctx, id)
}

func (repo *chatroomRepo) AddMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) error {
	return repo.store.Chatroom().AddMember(ctx, chatroom.Id, member)
}

func (repo *chatroomRepo) HasMember(ctx context.Context, chatroom *entity.Chatroom, member *entity.User) bool {
	return repo.store.Chatroom().HasMember(ctx, chatroom.Id, member)
}
