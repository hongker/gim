package internal

import (
	"flag"
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation"
	"gim/pkg/app"
	"gim/pkg/system"
	"log"
)

var (
	configFile = flag.String("conf", "./app.yaml", "configuration file")
)
func Run()  {
	flag.Parse()
	container := app.Container()

	infrastructure.Inject(container)
	application.Inject(container)
	presentation.Inject(container)

	if err := container.Invoke(func(conf *config.Config) error{
		return conf.LoadFile(*configFile)
	}); err != nil {
		panic(err)
	}

	err := container.Invoke(serve)
	if err != nil {
		panic(err)
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve(socket *presentation.Socket, conf *config.Config) error {
	return socket.Start(conf.Addr())
}
