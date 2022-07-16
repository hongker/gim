package app

import "go.uber.org/dig"

var container = NewContainer()

func Container() *dig.Container {
	return container
}
func NewContainer() *dig.Container {
	return dig.New()
}
