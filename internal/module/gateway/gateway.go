package gateway

import (
	"gim/internal/module/gateway/routes"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/server/grpc"
	"github.com/ebar-go/ego/server/http"
	"github.com/ebar-go/ego/server/ws"
)

type Instance struct {
	engine *ego.NamedEngine
	config *Config
}

func (instance *Instance) Start() {
	instance.prepare()

	instance.engine.Run()
}

func (instance *Instance) prepare() {
	// new http server
	httpServer := http.NewServer(instance.config.HttpServerAddress).
		RegisterRouteLoader(routes.Loader)

	// new grpc server
	grpcServer := grpc.NewServer(instance.config.GrpcServerAddress)

	// new socket server
	sockServer := ws.NewServer(instance.config.SockServerAddress)

	instance.engine.WithServer(httpServer, grpcServer, sockServer)
}
