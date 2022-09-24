package gateway

import (
	"gim/internal/module/gateway/http"
	"gim/internal/module/gateway/socket"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/utils/runtime"
)

type Instance struct {
	engine *ego.NamedEngine
	config *Config
}

func (instance *Instance) Start() {
	runtime.SetReallyCrash(false)
	instance.prepare()

	instance.engine.NonBlockingRun()
}

// initHttpServer initialize http server.
func (instance *Instance) initHttpServer() {
	// new http server
	httpServer := ego.NewHTTPServer(instance.config.HttpServerAddress).
		EnableCorsMiddleware().
		EnableTraceMiddleware(instance.config.TraceHeader).
		//EnableReleaseMode().
		EnableAvailableHealthCheck()

	if instance.config.EnablePprof {
		httpServer.EnablePprofHandler()
	}

	httpServer.RegisterRouteLoader(http.RouteLoader)

	instance.engine.WithServer(httpServer)
}

func (instance *Instance) initSockServer() {
	callback := socket.NewCallback()
	// new socket server
	sockServer := ego.NewWebsocketServer(instance.config.SockServerAddress).
		OnConnect(callback.OnConnect).
		OnDisconnect(callback.OnDisconnect).
		OnMessage(callback.OnMessage)

	instance.engine.WithServer(sockServer)
}

func (instance *Instance) prepare() {
	instance.initHttpServer()

	instance.initSockServer()
}
