package gateway

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
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
func (c *Controller) Run(stopCh <-chan struct{}, worker int) {
	c.once.Do(func() {
		c.initialize(worker)
	})
	c.run(stopCh)

}

// WithName set controller name.
func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

// initialize init controller dependencies.
func (c *Controller) initialize(worker int) {
	callback := NewCallback()

	wss := ego.NewWebsocketServer(c.config.Address).
		OnConnect(callback.OnConnect).
		OnDisconnect(callback.OnDisconnect).
		OnMessage(callback.OnMessage).WithWorker(worker)

	c.engine.WithServer(wss)
}

// run start module.
func (c *Controller) run(stopCh <-chan struct{}) {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)
	c.engine.NonBlockingRun()
	<-stopCh
	c.shutdown()
}

// shutdown shuts down the controller.
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}

// NewController returns a new Controller instance with *Config.
func NewController(config *Config) *Controller {
	return &Controller{
		name:   "gateway",
		config: config,
		engine: ego.New(),
	}
}
