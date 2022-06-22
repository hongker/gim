package infrastructure

import (
	"gim/internal/logic/infrastructure/persistence"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(persistence.NewUserRepository)
	_ = container.Provide(persistence.NewMessageRepo)
}
