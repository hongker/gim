package controller

// Controller represents module.
type Controller interface {
	Run(stopCh <-chan struct{})
}

// DaemonController run controller in daemon mode.
type DaemonController struct {
	Controller
}

func (c *DaemonController) NonBlockingRun() (stopCh chan struct{}) {
	stopCh = make(chan struct{})
	go c.Run(stopCh)
	return stopCh
}

func NewDaemonController(delegator Controller) *DaemonController {
	return &DaemonController{delegator}
}
