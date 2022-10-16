package framework

import (
	"gim/pkg/bytes"
)

// Engine represents context manager
type Engine struct {
	handleChains    []HandleFunc
	contextProvider ContextProvider
}

// Use registers middleware
func (e *Engine) Use(handler ...HandleFunc) {
	e.handleChains = append(e.handleChains, handler...)
}

// AcquireContext acquire context
func (e *Engine) AcquireContext() *Context {
	return e.contextProvider.AcquireContext()
}

// ------------------------private methods------------------------

func (e *Engine) processContext(ctx *Context) {
	// invoke handler chain
	e.handleChains[0](ctx)

	// release context after process
	e.releaseContext(ctx)
}

func (e *Engine) releaseContext(ctx *Context) {
	// release body
	bytes.Put(ctx.body)

	// release Context
	e.contextProvider.ReleaseContext(ctx)
}

func NewEngine() *Engine {
	engine := &Engine{}
	engine.contextProvider = NewSyncPoolContextProvider(func() interface{} {
		return &Context{engine: engine}
	})
	return engine
}
