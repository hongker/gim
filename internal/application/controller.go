package application

import (
	"gim/pkg/runtime/async"
)

// Controller is a controller manager for the core bootstrap application.
type Controller struct {
	runner *async.Runner

	// eventHandler is the event handler.
	eventHandler *EventHandler
}

func (c *Controller) Start() {
	if c.runner != nil {
		return
	}

	c.runner = async.NewRunner(c.eventHandler.Loop)

}

func (c *Controller) Stop() {
	c.runner.Stop()
}

type EventHandler struct{}

func (handler *EventHandler) Loop(stopCh chan struct{}) {

	<-stopCh
}
