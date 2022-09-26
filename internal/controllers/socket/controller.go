package socket

import (
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"sync"
)

type Controller struct {
	name string
	once sync.Once

	engine   *ego.NamedEngine
	callback *Callback
}

func (c *Controller) Run(stopCh <-chan struct{}, worker int) {
	c.once.Do(c.initialize)
	c.run()
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {
	wss := ego.NewWebsocketServer("").
		OnConnect(c.callback.OnConnect).
		OnDisconnect(c.callback.OnDisconnect).
		OnMessage(c.callback.OnMessage)

	c.engine.WithServer(wss)
}
func (c *Controller) run() {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)
	c.engine.NonBlockingRun()
}

func NewController() *Controller {
	return &Controller{
		name:     "default",
		engine:   ego.New(),
		callback: NewCallback(),
	}
}
