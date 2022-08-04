package event

import "sync"

// Event
type Event struct {
	// name
	Name string
	// event params
	Params interface{}
}

// Listener
type Listener struct {
	Mode    int
	Handler Handler
}

// Handler process event
type Handler func(params ...interface{})

// dispatcher
type dispatcher struct {
	items map[string][]Listener
	rmw   sync.RWMutex
}

var (
	// 初始化事件分发器，提前给map分配空间，减少因数组扩容带来的消耗
	instance = &dispatcher{items: make(map[string][]Listener, 100), rmw: sync.RWMutex{}}
)

// Register
func Register(eventName string, listener Listener) {
	instance.rmw.Lock()
	defer instance.rmw.Unlock()
	listeners, ok := instance.items[eventName]
	if !ok {
		// 预定义数组的长度为10
		listeners = make([]Listener, 0, 10)
	}
	listeners = append(listeners, listener)
	instance.items[eventName] = listeners
}

// Listen register a sync event
func Listen(eventName string, handler Handler) {
	Register(eventName, Listener{
		Handler: handler,
	})
}

// Has return event exist
func Has(eventName string) bool {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	_, ok := instance.items[eventName]
	return ok
}

// Trigger
func Trigger(eventName string, params ...interface{}) {
	instance.rmw.RLock()
	defer instance.rmw.RUnlock()
	listeners, ok := instance.items[eventName]
	if !ok {
		return
	}

	for _, listener := range listeners {
		listener.Handler(params...)

	}
}
