package gateway

import (
	"gim/framework"
	"github.com/ebar-go/ego/component"
	"sync"
)

// Controller represents gateway module.
type Controller struct {
	name string
	once sync.Once

	config *Config
}

// Run runs the controller.
func (c *Controller) Run(stopCh <-chan struct{}) {
	c.once.Do(c.initialize)

	component.Provider().Logger().Infof("controller running: [%s]", c.name)

	handler := NewHandler(c.config.HeartbeatInterval)
	app := framework.New(
		framework.WithConnectCallback(handler.OnConnect),
		framework.WithDisconnectCallback(handler.OnDisconnect),
	)

	handler.Install(app.Router())

	if err := app.Listen(framework.TCP, c.config.Address).Run(stopCh); err != nil {
		panic(err)
	}

	c.shutdown()

}

// WithName set controller name.
func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

// initialize init controller dependencies.
func (c *Controller) initialize() {

}

// shutdown shuts down the controller.
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
