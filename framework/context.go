package framework

import (
	"context"
	"gim/pkg/bytes"
	"log"
	"math"
	"sync"
)

type Context struct {
	context.Context
	container *ContextContainer
	conn      *Connection
	body      []byte
	index     int8
}

func (ctx *Context) Conn() *Connection {
	return ctx.conn
}

func (ctx *Context) Body() []byte {
	return ctx.body
}
func (ctx *Context) Run() {
	ctx.container.processContext(ctx)
}

func (ctx *Context) Next() {
	if ctx.index < maxIndex {
		ctx.index++
		ctx.container.handleChains[ctx.index](ctx)
	}
}
func (ctx *Context) Abort() {
	ctx.index = maxIndex
	log.Println("已被终止...")
}

func (ctx *Context) Reset(conn *Connection, body []byte) {
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

type ContextContainer struct {
	handleChains    []HandleFunc
	contextProvider ContextProvider
	packetMaxLength int
}

func (e *ContextContainer) Use(handler ...HandleFunc) {
	e.handleChains = append(e.handleChains, handler...)
}

func (e *ContextContainer) BuildContext(conn *Connection) (*Context, error) {
	buf := bytes.Get(e.packetMaxLength)
	n, err := conn.Read(buf)
	if err != nil {
		bytes.Put(buf)
		return nil, err
	}
	ctx := e.contextProvider.AcquireContext()
	ctx.Reset(conn, buf[:n])
	return ctx, nil
}

// ------------------------private methods------------------------

func (e *ContextContainer) processContext(ctx *Context) {
	e.handleChains[0](ctx)
	e.releaseContext(ctx)
}

func (e *ContextContainer) releaseContext(ctx *Context) {
	bytes.Put(ctx.body)
	e.contextProvider.ReleaseContext(ctx)
}

func NewContextContainer() *ContextContainer {
	engine := &ContextContainer{packetMaxLength: 512}
	engine.contextProvider = NewSyncPoolContextProvider(func() interface{} {
		return &Context{Context: context.Background(), container: engine}
	})
	return engine
}
