package controller

// Controller represents module.
type Controller interface {
	Run(stopCh <-chan struct{})
}
