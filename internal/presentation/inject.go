package presentation

import (
	"gim/internal/presentation/handler"
	"gim/internal/presentation/socket"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(socket.NewSocket)
	_ = container.Provide(handler.NewUserHandler)
	_ = container.Provide(handler.NewMessageHandler)
	_ = container.Provide(handler.NewGroupHandler)
}
