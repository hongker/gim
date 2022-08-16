package presentation

import (
	"gim/internal/presentation/handler"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(handler.NewUserHandler)
	_ = container.Provide(handler.NewMessageHandler)
	_ = container.Provide(handler.NewGroupHandler)
}
