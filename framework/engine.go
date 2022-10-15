package framework

import (
	"context"
	"errors"
	"github.com/ebar-go/ego/utils/runtime"
	"log"
)

// Engine represents im framework public access api.
type Engine struct {
	schemas  []Schema
	callback *Callback
	codec    Codec
	router   *Router
	event    *Event
}

// WithProtocol set different protocol
func (engine *Engine) Listen(protocol Protocol, addr string) *Engine {
	engine.schemas = append(engine.schemas, NewSchema(protocol, addr))
	return engine
}

// WithCallback use callback
func (engine *Engine) WithCallback(callback *Callback) *Engine {
	engine.callback = callback
	return engine
}

// WithCodec use codec to pack/unpack message.
func (engine *Engine) WithCodec(codec Codec) *Engine {
	engine.codec = codec
	return engine
}

// WithRouter set router
func (engine *Engine) WithRouter(router *Router) *Engine {
	engine.router = router
	return engine
}

// WithEvent set event
func (engine *Engine) WithEvent(event *Event) *Engine {
	engine.event = event
	return engine
}

// Start starts the engine
func (engine *Engine) Run(stopCh <-chan struct{}) error {
	if len(engine.schemas) == 0 {
		return errors.New("empty listen target")
	}

	// listen protocol
	schemaContext, schemeCancel := context.WithCancel(context.Background())
	for _, schema := range engine.schemas {
		err := schema.Listen(schemaContext.Done())
		runtime.HandleError(err, func(err error) {
			log.Println("listen error:", err)
		})
	}

	log.Println("engine started")
	runtime.WaitClose(stopCh, schemeCancel, engine.Stop)

	return nil
}

// Stop shuts down the engine.
func (engine *Engine) Stop() {}

// New returns a new engine instance
func New() *Engine {
	return &Engine{}
}
