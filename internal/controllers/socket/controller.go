package socket

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"sync"
)

type Controller struct {
	name string
	once sync.Once

	config *Config

	engine   *ego.NamedEngine
	callback *Callback
}

func (c *Controller) Run(stopCh <-chan struct{}, worker int) {
	c.once.Do(c.initialize)
	c.run(stopCh)
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {
	wss := ego.NewWebsocketServer(c.config.Address).
		OnConnect(c.callback.OnConnect).
		OnDisconnect(c.callback.OnDisconnect).
		OnMessage(c.callback.OnMessage)

	c.engine.WithServer(wss)
}
func (c *Controller) run(stopCh <-chan struct{}) {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)
	c.engine.NonBlockingRun()
	<-stopCh
}

func NewController(config *Config) *Controller {
	return &Controller{
		name:     "default",
		config:   config,
		engine:   ego.New(),
		callback: NewCallback(),
	}
}
