package framework

import "errors"

// Engine represents im framework public access api.
type Engine struct {
	schemas []Schema
}

// WithProtocol set different protocol
func (engine *Engine) Listen(protocol Protocol, addr string) *Engine {
	engine.schemas = append(engine.schemas, NewSchema(protocol, addr))
	return engine
}

// WithCallback use callback
func (engine *Engine) WithCallback(callback *Callback) *Engine { return engine }

// WithCodec use codec to pack/unpack message.
func (engine *Engine) WithCodec(codec Codec) *Engine { return engine }

// WithRouter set router
func (engine *Engine) WithRouter(router *Router) *Engine { return engine }

// WithEvent set event
func (engine *Engine) WithEvent(event *Event) *Engine { return engine }

// Start starts the engine
func (engine *Engine) Run(stopCh <-chan struct{}) error {
	if len(engine.schemas) == 0 {
		return errors.New("empty listen target")
	}

	<-stopCh
	engine.Stop()
	return nil
}

// Stop shuts down the engine.
func (engine *Engine) Stop() {}

// New returns a new engine instance
func New() *Engine {
	return &Engine{}
}
