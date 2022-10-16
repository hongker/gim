package framework

import (
	"context"
	"log"
	"math"
	"sync"
)

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(p []byte, v interface{}) error
}

type Context struct {
	context.Context
	engine *Engine
	conn   *Connection
	body   []byte
	index  int8
}

func (ctx *Context) Conn() *Connection {
	return ctx.conn
}

func (ctx *Context) Body() []byte {
	return ctx.body
}
func (ctx *Context) Run() {
	ctx.engine.processContext(ctx)
}

func (ctx *Context) Next() {
	if ctx.index < maxIndex {
		ctx.index++
		ctx.engine.handleChains[ctx.index](ctx)
	}
}
func (ctx *Context) Abort() {
	ctx.index = maxIndex
	log.Println("已被终止...")
}

func (ctx *Context) reset(conn *Connection, body []byte) {
	ctx.index = 0
	ctx.body = body
	ctx.conn = conn
	ctx.Context = context.Background()
}

const (
	maxIndex = math.MaxInt8 / 2
)

type HandleFunc func(ctx *Context)

type ContextProvider interface {
	AcquireContext() *Context
	ReleaseContext(ctx *Context)
}

type SyncPoolContextProvider struct {
	pool *sync.Pool
}

func (provider *SyncPoolContextProvider) AcquireContext() *Context {
	return provider.pool.Get().(*Context)
}

func (provider *SyncPoolContextProvider) ReleaseContext(ctx *Context) {
	provider.pool.Put(ctx)
}

func NewSyncPoolContextProvider(constructor func() interface{}) ContextProvider {
	return &SyncPoolContextProvider{pool: &sync.Pool{New: constructor}}
}
