package interfaces

import (
	"gim/internal/interfaces/handler"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(NewSocket)
	_ = container.Provide(handler.NewUserHandler)
	_ = container.Provide(handler.NewMessageHandler)
}
