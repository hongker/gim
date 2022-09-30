package gateway

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

// Controller represents gateway module.
type Controller struct {
	name string
	once sync.Once

	config *Config

	engine *ego.NamedEngine
}

// Run runs the controller.
func (c *Controller) Run(stopCh <-chan struct{}) {
	c.once.Do(c.initialize)
	c.run()

	runtime.WaitClose(stopCh, c.shutdown)

}

// WithName set controller name.
func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

// initialize init controller dependencies.
func (c *Controller) initialize() {
	c.engine = ego.New()

	callback := NewCallback(c.config.HeartbeatInterval)

	wss := ego.NewWebsocketServer(c.config.Address).
		WithWorker(c.config.WorkerNumber).
		OnConnect(callback.OnConnect).
		OnDisconnect(callback.OnDisconnect).
		OnMessage(callback.OnMessage)

	c.engine.WithServer(wss)
}

// run start module.
func (c *Controller) run() {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)
	c.engine.NonBlockingRun()
}

// shutdown shuts down the controller.
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
