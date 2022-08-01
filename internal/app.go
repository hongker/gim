package internal

import (
	"flag"
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

	err := container.Invoke(serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve() error  {
	socket := interfaces.NewSocket(*addr)

	return utils.Execute(socket.Start)
}
