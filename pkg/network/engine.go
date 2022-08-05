package network

import (
	"math"
	"math/rand"
	"sync"
)

const maxIndex = math.MaxInt8 / 2

type Engine struct {
	handleChains []HandleFunc
	ctxPools []*sync.Pool
}

func newEngine(poolSize int) *Engine {
	e := &Engine{
		handleChains: make([]HandleFunc, 0, 10),
		ctxPools: make([]*sync.Pool, poolSize),
	}
	e.init()
	return e
}

func (e *Engine) init() {
	for i := 0; i < len(e.ctxPools); i++ {
		e.ctxPools[i] = &sync.Pool{New: func() interface{} {
			return e.allocateContext()
		}}
	}
}

func (e *Engine) Use(handler ...HandleFunc) {
	e.handleChains = append(e.handleChains, handler...)
}

func (e *Engine) allocateContext() *Context {
	return &Context{engine: e}
}

func (e *Engine) contextPool() *sync.Pool {
	return e.ctxPools[rand.Intn(len(e.ctxPools)-1)]
}
