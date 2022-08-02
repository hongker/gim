package internal

import (
	"flag"
	"gim/internal/applications"
	"gim/internal/infrastructure"
	"gim/internal/interfaces"
	"gim/pkg/app"
	"gim/pkg/system"
	"gim/pkg/utils"
	"log"
)

var (
	addr = flag.String("addr", "0.0.0.0:8088", "socket address")
)
func Run()  {
	flag.Parse()
	container := app.Container()

	infrastructure.Inject(container)
	applications.Inject(container)
	interfaces.Inject(container)

	err := container.Invoke(serve)
	if err != nil {
		panic(err)
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve(socket *interfaces.Socket) error  {
	return utils.Execute(func() error {
		return socket.Start(*addr)
	})
}
