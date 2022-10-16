package job

import (
	"gim/internal/controller/job/task"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
)

// Controller represents cronjob controller.
type Controller struct {
	name string

	config *Config
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)

	task.NewMessageTask(c.config.QueuePollInterval, c.config.QueuePollCount).Start()

	runtime.WaitClose(stopCh, c.shutdown)
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
