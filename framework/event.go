package framework

type Event struct{}

func (e *Event) Listen(name string)  {}
func (e *Event) Trigger(name string) {}
