package event

import "sync"

var (
	defaultInstance = NewDispatcher()
	Listen          = defaultInstance.Listen
	Trigger         = defaultInstance.Trigger
)

// Handler process event
type Handler func(param any)
type Dispatcher struct {
	items map[string][]Handler
	rmw   sync.RWMutex
}

// Listen register a sync event
func (instance *Dispatcher) Listen(eventName string, handler Handler) {
	instance.rmw.Lock()
	defer instance.rmw.Unlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		// 预定义数组的长度为10
		handlers = make([]Handler, 0, 10)
	}
	handlers = append(handlers, handler)
	instance.items[eventName] = handlers
}

// Has return event exist
func (instance *Dispatcher) Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger make event trigger with given name and params
func (instance *Dispatcher) Trigger(eventName string, param any) {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	handlers, ok := instance.items[eventName]
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(param)
	}
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{items: make(map[string][]Handler)}
}

type Event[T any] struct {
	Name string
}

// NewEvent creates a new Event with the given name.
func NewEvent[T any](name string) Event[T] {
	return Event[T]{Name: name}
}

// Bind binds handler with default dispatcher
func (e Event[T]) Bind(handler func(param T)) {
	e.BindWithDispatcher(handler, defaultInstance)
}

func (e Event[T]) BindWithDispatcher(handler func(param T), dispatcher *Dispatcher) {
	dispatcher.Listen(e.Name, func(param any) {
		data, ok := param.(T)
		if !ok {
			return
		}
		handler(data)
	})
}
