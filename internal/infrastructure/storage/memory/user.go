package memory

import (
	"context"
	"gim/internal/domain/entity"
	"github.com/ebar-go/ego/utils/structure"
)

type UserStorage struct {
	container *structure.ConcurrentMap[*entity.User]
}

func NewUserStorage() *UserStorage {
	return &UserStorage{container: structure.NewConcurrentMap[*entity.User]()}
}

func (storage *UserStorage) Create(ctx context.Context, item *entity.User) error {
	return storage.container.Set(item.Id, item)
}

func (storage *UserStorage) Find(ctx context.Context, id string) (*entity.User, error) {
	return storage.container.Find(id)
}
