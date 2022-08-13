package internal

import (
	"flag"
	"gim/internal/application"
	"gim/internal/infrastructure"
	"gim/internal/infrastructure/config"
	"gim/internal/presentation"
	"gim/pkg/app"
	"gim/pkg/system"
	"gim/pkg/utils"
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

	err := container.Invoke(serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}


func serve(socket *presentation.Socket, conf *config.Config) error {
	return utils.Execute(func() error {
		return conf.LoadFile(*configFile)
	}, socket.Start)

}
