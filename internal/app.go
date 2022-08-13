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

	port = flag.Int("port", 8080, "server port")
	maxLimit = flag.Int("max-limit", 1000, "max number of session history messages")
)

var (
	Version = "1.0.0"
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
	log.Printf("version: %v is running...", Version)
	if err := conf.LoadFile(*configFile); err != nil {
		return err
	}

	conf.Server.Port = *port
	conf.Message.MaxStoreSize = *maxLimit

	return socket.Start()

}
