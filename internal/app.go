package internal

import (
	"flag"
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/presentation"
	"gim/pkg/app"
	"gim/pkg/system"
	"log"
)

var (
	addr = flag.String("addr", "0.0.0.0:8088", "socket address")
	pluginStore = flag.String("plugin-store", "memory", "plugin store")
)
func Run()  {
	flag.Parse()
	container := app.Container()


	infrastructure.Inject(container, *pluginStore)
	application.Inject(container)
	presentation.Inject(container)

	err := container.Invoke(func(socket *presentation.Socket) error {
		return socket.Start(*addr)
	})

	if err != nil {
		panic(err)
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

