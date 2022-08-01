package interfaces

import "go.uber.org/dig"

func Inject(container *dig.Container)  {
	_ = container.Provide(NewSocket)
}
