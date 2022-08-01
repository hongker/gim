package infrastructure

import (
	"gim/internal/infrastructure/cache"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(cache.NewMessageRepo)
	_ = container.Provide(cache.NewUserRepo)
}
