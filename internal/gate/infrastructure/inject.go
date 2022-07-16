package infrastructure

import (
	"gim/internal/gate/infrastructure/config"
	"gim/internal/gate/infrastructure/grpc"
	"go.uber.org/dig"
)

func Inject(container *dig.Container)  {
	_ = container.Provide(config.NewConfig)
	_ = container.Provide(grpc.NewLogicClient)
}
