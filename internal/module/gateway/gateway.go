package gateway

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/server/http"
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
	httpServer := http.NewServer(instance.config.HttpServerAddress)

	instance.engine.WithServer(httpServer)
}
