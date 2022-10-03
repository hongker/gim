package memory

import (
	"context"
	"gim/internal/domain/entity"
	"gim/pkg/store"
	"github.com/ebar-go/ego/utils/convert"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/structure"
)

type MessageStorage struct {
	container *structure.ConcurrentMap[*entity.Message]
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{container: structure.NewConcurrentMap[*entity.Message]()}
}

func (storage *MessageStorage) Create(ctx context.Context, msg *entity.Message) error {
	return storage.container.Set(msg.Id, msg)
}

func (storage *MessageStorage) Find(ctx context.Context, id string) (*entity.Message, error) {
	return storage.container.Find(id)
}

type SessionMessage struct {
	session  *entity.Session
	messages []string
}

type UserSession struct {
	store.Set
}

type SessionStorage struct {
	maxSessionMessageSize int
	container             *structure.ConcurrentMap[*SessionMessage]
	userContainer         *structure.ConcurrentMap[*UserSession]
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		maxSessionMessageSize: 100,
		container:             structure.NewConcurrentMap[*SessionMessage](),
		userContainer:         structure.NewConcurrentMap[*UserSession](),
	}
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
			userSession = &UserSession{store.ThreadSafe(item.Id)}
			return storage.userContainer.Set(uid, userSession)
		}

		userSession.Add(item.Id)
		return nil
	})
}

func (storage *SessionStorage) List(ctx context.Context, uid string) ([]*entity.Session, error) {
	userSession, exist := storage.userContainer.Get(uid)
	emptySessions := make([]*entity.Session, 0)
	if !exist {
		return emptySessions, nil
	}

	sessionIds := convert.ToSlice[string](userSession.ToSlice())

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
