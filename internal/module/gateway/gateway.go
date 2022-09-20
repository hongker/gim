package gateway

import (
	"gim/internal/module/gateway/handler"
	"gim/internal/module/gateway/route"
	"gim/internal/module/gateway/socket"
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

	instance.engine.NonBlockingRun()
}

// initHttpServer initialize http server.
func (instance *Instance) initHttpServer() {
	// register handlers
	route.Container().RegisterHandler(
		handler.NewUserHandler(),
		handler.NewMessageHandler(),
		handler.NewChatRoomHandler(),
	)

	// new http server
	httpServer := http.NewServer(instance.config.HttpServerAddress).
		EnableCorsMiddleware().
		EnableTraceMiddleware(instance.config.TraceHeader).
		//EnableReleaseMode().
		EnableAvailableHealthCheck()

	if instance.config.EnablePprof {
		httpServer.EnablePprofHandler()
	}
	httpServer.RegisterRouteLoader(route.Loader)

	instance.engine.WithServer(httpServer)
}

func (instance *Instance) initSockServer() {
	callback := socket.NewCallback()
	// new socket server
	sockServer := ws.NewServer(instance.config.SockServerAddress).
		OnConnect(callback.OnConnect).
		OnDisconnect(callback.OnDisconnect).
		OnMessage(callback.OnMessage)

	instance.engine.WithServer(sockServer)
}

func (instance *Instance) prepare() {
	instance.initHttpServer()

	instance.initSockServer()
}
