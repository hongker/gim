package main

import (
	"context"
	"gim/framework"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"log"
)

func main() {
	app := framework.New(framework.WithConnectCallback(func(conn *framework.Connection) {
		log.Printf("[%s] Connected\n", conn.UUID())
	}), framework.WithDisconnectCallback(func(conn *framework.Connection) {
		log.Printf("[%s] Disconnected\n", conn.UUID())
	}))

	NewController().Install(app.Router())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := app.Listen(framework.TCP, ":8090").
			Listen(framework.WEBSOCKET, ":8091").
			Run(ctx.Done()); err != nil {
			panic(err)
		}
	}()

	runtime.WaitClose(signal.SetupSignalHandler())

}
