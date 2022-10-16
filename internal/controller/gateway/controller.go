package gateway

import (
	"gim/framework"
	"github.com/ebar-go/ego/component"
)

// Controller represents gateway module.
type Controller struct {
	name string

	config *Config
}

// Run runs the controller.
func (c *Controller) Run(stopCh <-chan struct{}) {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)

	handler := NewHandler(c.config.HeartbeatInterval)
	app := framework.New(
		framework.WithConnectCallback(handler.OnConnect),
		framework.WithDisconnectCallback(handler.OnDisconnect),
		framework.WithMiddleware(handler.checkLogin),
	)

	handler.Install(app.Router())

	if c.config.TCPAddress != "" {
		app.Listen(framework.TCP, c.config.TCPAddress)
	}

	if c.config.WebsocketAddress != "" {
		app.Listen(framework.WEBSOCKET, c.config.WebsocketAddress)
	}
	if err := app.Run(stopCh); err != nil {
		panic(err)
	}

	c.shutdown()

}

// WithName set controller name.
func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

// shutdown shuts down the controller.
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
