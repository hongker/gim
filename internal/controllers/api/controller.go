package api

import (
	"github.com/ebar-go/ego/component"
	"sync"
)

type Controller struct {
	name string
	once sync.Once
}

func (c *Controller) Run(stopCh <-chan struct{}, worker int) {
	c.once.Do(c.initialize)
	c.run()
	<-stopCh
	c.shutdown()

}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {}
func (c *Controller) run() {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)
}
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}

func NewController() *Controller {
	return &Controller{name: "api"}
}
