package main

import (
	"context"
	"gim/framework"
	"gim/framework/codec"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/ebar-go/ego/utils/runtime/signal"
	uuid "github.com/satori/go.uuid"
	"log"
)

func main() {
	app := framework.New(framework.WithConnectCallback(func(conn *framework.Connection) {
		log.Printf("[%s] Connected\n", conn.UUID())
	}), framework.WithDisconnectCallback(func(conn *framework.Connection) {
		log.Printf("[%s] Disconnected\n", conn.UUID())
	}))

	func(router *framework.Router) {
		router.WithCodec(codec.Default()).OnNotFound(func(ctx *framework.Context) {
			log.Println("operation not found")
		}).OnError(func(ctx *framework.Context, err error) {
			log.Println("operation error: ", ctx.Operate(), err)
		})

		router.Route(1, framework.StandardHandler[LoginRequest, LoginResponse](Login))

	}(app.Router())

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

type LoginRequest struct{ Name string }
type LoginResponse struct {
	ID    string
	Token string
}

func Login(ctx *framework.Context, req *LoginRequest) (*LoginResponse, error) {
	return &LoginResponse{ID: "1001", Token: uuid.NewV4().String()}, nil
}
