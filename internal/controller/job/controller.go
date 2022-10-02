package job

import (
	"gim/internal/controller/job/task"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"sync"
)

// Controller represents cronjob controller.
type Controller struct {
	name string
	once sync.Once

	config *Config
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	c.once.Do(c.initialize)
	c.run()

	runtime.WaitClose(stopCh, c.shutdown)
}

func (c *Controller) WithName(name string) *Controller {
	c.name = name
	return c
}

func (c *Controller) initialize() {
	task.NewMessageTask(c.config.QueuePollInterval, c.config.QueuePollCount).Start()
}

func (c *Controller) run() {
	component.Provider().Logger().Infof("controller running: [%s]", c.name)

}
func (c *Controller) shutdown() {
	component.Provider().Logger().Infof("controller shutdown: [%s]", c.name)
}
