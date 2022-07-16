package gate

// 长链接网关，负责维护客户端的连接、数据收发
// 收到客户端请求后，通过grpc调用，将请求数据发送给logic服务，并将响应数据返回给客户端

import (
	"gim/internal/gate/applications"
	"gim/internal/gate/infrastructure"
	"gim/internal/gate/interfaces"
	"gim/pkg/app"
	"gim/pkg/errgroup"
	"gim/pkg/system"
	"log"
)

func Run() {
	container := app.Container()

	infrastructure.Inject(container)
	applications.Inject(container)
	interfaces.Inject(container)

	err := container.Invoke(serve)
	system.SecurePanic(err)

	system.Shutdown(func() {
		log.Println("server shutdown")
	})
}

func serve(socket *interfaces.Socket, server *interfaces.GRPCServer) error {
	return errgroup.Do(socket.Start, server.Start)
}