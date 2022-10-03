package memory

import (
	"context"
	"gim/internal/domain/entity"
	"gim/pkg/store"
	"github.com/ebar-go/ego/utils/structure"
)

type Chatroom struct {
	chatroom *entity.Chatroom
	members  store.Set
}

type ChatroomStorage struct {
	container *structure.ConcurrentMap[*Chatroom]
}

func NewChatroomStorage() *ChatroomStorage {
	return &ChatroomStorage{container: structure.NewConcurrentMap[*Chatroom]()}
}

func (storage *ChatroomStorage) Create(ctx context.Context, item *entity.Chatroom) error {
	return storage.container.Set(item.Id, &Chatroom{chatroom: item, members: store.ThreadSafe(item.Creator)})
}

func (storage *ChatroomStorage) Find(ctx context.Context, id string) (*entity.Chatroom, error) {
	item, err := storage.container.Find(id)
	if err != nil {
		return nil, err
	}
	return item.chatroom, nil
}

func (storage *ChatroomStorage) GetMember(ctx context.Context, id string) ([]string, error) {
	item, err := storage.container.Find(id)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, item.members.Size())
	members := item.members.ToSlice()
	for _, member := range members {
		res = append(res, member.(string))
	}
	return res, nil
}

func (storage *ChatroomStorage) AddMember(ctx context.Context, id string, member *entity.User) error {
	item, err := storage.container.Find(id)
	if err != nil {
		return err
	}
	item.members.Add(member.Id)
	return nil
}

func (storage *ChatroomStorage) RemoveMember(ctx context.Context, id string, member *entity.User) error {
	item, err := storage.container.Find(id)
	if err != nil {
		return err
	}
	item.members.Remove(member.Id)
	return nil
}

func (storage *ChatroomStorage) HasMember(ctx context.Context, id string, member *entity.User) bool {
	item, err := storage.container.Find(id)
	if err != nil {
		return false
	}
	return item.members.Contain(member.Id)
}
