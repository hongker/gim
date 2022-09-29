package memory

import (
	"context"
	"gim/internal/domain/entity"
	"gim/pkg/store"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

type MapContainer[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func (container *MapContainer[T]) Get(key string) (T, bool) {
	container.mu.RLock()
	defer container.mu.RUnlock()
	item, ok := container.items[key]
	return item, ok
}

func (container *MapContainer[T]) Set(key string, value T) error {
	container.mu.Lock()
	container.items[key] = value
	container.mu.Unlock()
	return nil
}

func (container *MapContainer[T]) Find(key string) (T, error) {
	item, exist := container.Get(key)
	var empty T
	if !exist {
		return empty, errors.NotFound("not found")
	}
	return item, nil
}

func NewMapContainer[T any]() *MapContainer[T] {
	return &MapContainer[T]{items: make(map[string]T, 0)}
}

type SetContainer[T string | int] struct {
	set store.Set
}

func (container *SetContainer[T]) Add(item T) {
	container.set.Add(item)
}
func (container *SetContainer[T]) Remove(item T) {
	container.set.Remove(item)
}
func (container *SetContainer[T]) Contain(item T) bool {
	return container.set.Contain(item)
}

func ToSlice[T any](items []interface{}) []T {
	res := make([]T, 0)
	for _, item := range items {
		res = append(res, item.(T))
	}
	return res
}

type UserStorage struct {
	container *MapContainer[*entity.User]
}

func NewUserStorage() *UserStorage {
	return &UserStorage{container: NewMapContainer[*entity.User]()}
}

func (storage *UserStorage) Create(ctx context.Context, item *entity.User) error {
	return storage.container.Set(item.Id, item)
}

func (storage *UserStorage) Find(ctx context.Context, id string) (*entity.User, error) {
	return storage.container.Find(id)
}

type MessageStorage struct {
	container *MapContainer[*entity.Message]
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{container: NewMapContainer[*entity.Message]()}
}

func (storage *MessageStorage) Create(ctx context.Context, msg *entity.Message) error {
	return storage.container.Set(msg.Id, msg)
}

func (storage *MessageStorage) Find(ctx context.Context, id string) (*entity.Message, error) {
	return storage.container.Find(id)
}

type Chatroom struct {
	chatroom *entity.Chatroom
	members  store.Set
}
type ChatroomStorage struct {
	container *MapContainer[*Chatroom]
}

func NewChatroomStorage() *ChatroomStorage {
	return &ChatroomStorage{container: NewMapContainer[*Chatroom]()}
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

type SessionMessage struct {
	session  *entity.Session
	messages []string
}
type SessionStorage struct {
	maxSessionMessageSize int
	container             *MapContainer[*SessionMessage]
	userContainer         *MapContainer[*SetContainer[string]]
}

func (storage *SessionStorage) Create(ctx context.Context, uid string, item *entity.Session) error {
	return runtime.Call(func() error {
		_, exist := storage.container.Get(item.Id)
		if exist {
			return nil
		}
		return storage.container.Set(item.Id, &SessionMessage{session: item, messages: make([]string, 0, 100)})
	}, func() error {
		userSession, exist := storage.userContainer.Get(uid)
		if !exist {
			userSession = &SetContainer[string]{set: store.ThreadSafe(item.Id)}
			return storage.userContainer.Set(uid, userSession)
		}

		userSession.Add(item.Id)
		return nil
	})
}

func (storage *SessionStorage) List(ctx context.Context, uid string) ([]*entity.Session, error) {
	container, exist := storage.userContainer.Get(uid)
	emptySessions := make([]*entity.Session, 0)
	if !exist {
		return emptySessions, nil
	}

	sessionIds := ToSlice[string](container.set.ToSlice())

	for _, id := range sessionIds {
		session, err := storage.container.Find(id)
		if err != nil {
			continue
		}
		if len(session.messages) > 0 {
			session.session.Last = session.messages[len(session.messages)-1]
		}

		emptySessions = append(emptySessions, session.session)
	}
	return emptySessions, nil

}

func (storage *SessionStorage) SaveMessage(ctx context.Context, sessionId, msgId string) error {
	item, err := storage.container.Find(sessionId)
	if err != nil {
		return err
	}

	item.messages = append(item.messages, msgId)
	if overflow := len(item.messages) - storage.maxSessionMessageSize; overflow > 0 {
		item.messages = item.messages[overflow:]
	}
	return nil

}

func (storage *SessionStorage) QueryMessage(ctx context.Context, sessionId string) ([]string, error) {
	item, exist := storage.container.Get(sessionId)
	emptyMessageIds := make([]string, 0)
	if !exist {
		return emptyMessageIds, nil
	}
	return item.messages, nil
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		maxSessionMessageSize: 100,
		container:             NewMapContainer[*SessionMessage](),
		userContainer:         NewMapContainer[*SetContainer[string]](),
	}
}
