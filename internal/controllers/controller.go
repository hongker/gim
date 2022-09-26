package controllers

// Controller represents module.
type Controller interface {
	Run(stopCh <-chan struct{}, worker int)
}
