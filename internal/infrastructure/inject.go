package infrastructure

import (
	"gim/internal/infrastructure/memory"
	"go.uber.org/dig"
)

func Inject(container *dig.Container, storeType string)  {
	_ = container.Provide(memory.NewMessageRepo)
	_ = container.Provide(memory.NewUserRepo)
	_ = container.Provide(memory.NewGroupRepo)
	_ = container.Provide(memory.NewGroupUserRepo)
}
