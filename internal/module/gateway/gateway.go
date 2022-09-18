package gateway

import (
	"gim/internal/module/gateway/handler"
	"gim/internal/module/gateway/route"
	"github.com/ebar-go/ego"
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

// initHttpServer initialize http server.
func (instance *Instance) initHttpServer() {
	// register handlers
	route.Container().RegisterHandler(handler.NewUserHandler())

	// new http server
	httpServer := http.NewServer(instance.config.HttpServerAddress).
		EnableCorsMiddleware().
		EnableTraceMiddleware("trace").
		//EnableReleaseMode().
		EnablePprofHandler().
		EnableAvailableHealthCheck().
		RegisterRouteLoader(route.Loader)

	instance.engine.WithServer(httpServer)
}

func (instance *Instance) initSockServer() {
	// new socket server
	sockServer := ws.NewServer(instance.config.SockServerAddress)

	instance.engine.WithServer(sockServer)
}

func (instance *Instance) prepare() {
	instance.initHttpServer()

	instance.initSockServer()
}
