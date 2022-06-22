package application

import "go.uber.org/dig"

func Inject(container *dig.Container) {
	_ = container.Provide(NewMessage)
	_ = container.Provide(NewUserApp)
}
