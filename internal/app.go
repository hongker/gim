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
	addr = flag.String("addr", "0.0.0.0:8088", "socket address")
	messagePushCount = flag.Int("msg-push-count", 10, "number of message pushed")
	messageStoreCount = flag.Int("msg-store-count", 10000, "number of message stored")
)
func Run()  {
	flag.Parse()
	container := app.Container()

	infrastructure.Inject(container)
	application.Inject(container)
	presentation.Inject(container)

	err := container.Invoke(serve)

	if err != nil {
		panic(err)
	}

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve(socket *presentation.Socket, conf *config.Config) error {
	options := []config.Option{}
	if *messagePushCount > 0 {
		options = append(options, config.WithMessagePushCount(*messageStoreCount))
	}

	if *messageStoreCount > 0 {
		options = append(options, config.WithMessageMaxStoreSize(*messageStoreCount))
	}
	conf.WithOptions(options...)

	return socket.Start(*addr)
}
