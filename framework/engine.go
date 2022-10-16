package framework

import (
	"context"
	"gim/pkg/bytes"
)

// Engine represents context manager
type Engine struct {
	packetMaxLength int
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
	e.handleChains[0](ctx)
	e.releaseContext(ctx)
}

func (e *Engine) releaseContext(ctx *Context) {
	bytes.Put(ctx.body)
	e.contextProvider.ReleaseContext(ctx)
}

func NewEngine() *Engine {
	engine := &Engine{packetMaxLength: 512}
	engine.contextProvider = NewSyncPoolContextProvider(func() interface{} {
		return &Context{Context: context.Background(), engine: engine}
	})
	return engine
}
