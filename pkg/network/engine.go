package network

import (
	"math"
	"sync"
)

const maxIndex = math.MaxInt8 / 2

type Engine struct {
	handleChains []HandleFunc
}

func (e *Engine) Use(handler ...HandleFunc) {
	e.handleChains = append(e.handleChains, handler...)
}

func (e *Engine) allocateContext() *Context {
	return &Context{engine: e}
}

func (e *Engine) ContextPool() sync.Pool {
	return sync.Pool{New: func() interface{} {
		return e.allocateContext()
	}}
}
