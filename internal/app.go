package internal

import (
	"flag"
	"gim/internal/aggregate"
	"gim/internal/infrastructure"
	"gim/internal/interfaces"
	"gim/pkg/app"
	"gim/pkg/system"
	"log"
)

var (
	addr = flag.String("addr", "0.0.0.0:8088", "socket address")
)
func Run()  {
	flag.Parse()
	container := app.Container()

	infrastructure.Inject(container)
	aggregate.Inject(container)
	interfaces.Inject(container)

	err := container.Invoke(func(socket *interfaces.Socket) error {
		return socket.Start(*addr)
	})

	if err != nil {
		panic(err)
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

