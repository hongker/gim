package internal

import (
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation"
	"gim/internal/presentation/socket"
	"gim/pkg/app"
	"gim/pkg/system"
)

func Run(conf *config.Config) error {
	if conf.Debug {
		go system.ShowMemoryUsage()
	}
	config.Initialize(conf)

	container := app.Container()

	infrastructure.Inject(container)
	application.Inject(container)
	presentation.Inject(container)

	if err := container.Invoke(socket.Initialize); err != nil {
		return err
	}

	return socket.Get().Start()
}
