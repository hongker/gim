package infrastructure

import (
	"gim/internal/infrastructure/config"
	"gim/internal/infrastructure/memory"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(config.New)
	_ = container.Provide(memory.NewMessageRepo)
	_ = container.Provide(memory.NewUserRepo)
	_ = container.Provide(memory.NewGroupRepo)
	_ = container.Provide(memory.NewGroupUserRepo)
}
