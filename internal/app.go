package internal

import (
	"gim/internal/interfaces"
	"gim/pkg/app"
	"gim/pkg/errgroup"
	"gim/pkg/system"
	"log"
)

func Run()  {
	container := app.Container()

	err := container.Invoke(serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve() error  {
	socket := interfaces.NewSocket("")

	return errgroup.Do(socket.Start)
}
