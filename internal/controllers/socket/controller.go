package socket

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
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {}
func (c *Controller) run() {
	component.Provider().Logger().Infof("socket controller [%s] running", c.name)
}

func NewController() *Controller {
	return &Controller{name: "default"}
}
