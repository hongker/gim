package aggregate

import "go.uber.org/dig"

func Inject(container *dig.Container)  {
	_ = container.Provide(NewUserApp)
	_ = container.Provide(NewMessageApp)
	_ = container.Provide(NewGateApp)
	_ = container.Provide(NewGroupApp)
}
