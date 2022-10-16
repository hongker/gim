package framework

import (
	"context"
	"gim/framework/codec"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type Handler func(ctx *Context) (any, error)

func StandardHandler[Request, Response any](action func(ctx context.Context, request *Request) (*Response, error)) Handler {
	return func(ctx *Context) (any, error) {
		request := new(Request)
		if err := ctx.Bind(request); err != nil {
			return nil, err
		}
		return action(ctx, request)

	}
}

// Router
type Router struct {
	rwm          sync.RWMutex
	handlers     map[int16]Handler
	codec        codec.Codec
	errorHandler func(ctx *Context, err error)
}

// Route register handler for operate
func (router *Router) Route(operate int16, handler Handler) *Router {
	router.rwm.Lock()
	router.handlers[operate] = handler
	router.rwm.Unlock()
	return router
}

func (router *Router) Request(ctx *Context) {
	// unpack
	packet, err := router.codec.Unpack(ctx.body)
	if err != nil {
		router.handleError(ctx, errors.WithMessage(err, "invalid request"))
		return
	}

	// match handler
	router.rwm.RLock()
	handler, ok := router.handlers[packet.Operate]
	router.rwm.RUnlock()
	if !ok {
		router.handleError(ctx, errors.Errorf("operation not allowed:%v", packet.Operate))
		return
	}

	ctx.packet = packet
	response, err := handler(ctx)
	if err != nil {
		router.handleError(ctx, errors.WithMessage(err, "handle operation"))
		return
	}

	packet.Operate++
	packet.Seq++
	// pack response
	msg, err := router.codec.Pack(packet, response)
	if err != nil {
		router.handleError(ctx, errors.WithMessage(err, "invalid response"))
		return
	}
	ctx.Conn().Push(msg)
}

func (router *Router) OnError(handler func(ctx *Context, err error)) *Router {
	router.errorHandler = handler
	return router
}

func (router *Router) handleError(ctx *Context, err error) {
	if router.errorHandler != nil {
		router.errorHandler(ctx, err)
	}
}

func NewRouter() *Router {
	return &Router{
		handlers: map[int16]Handler{},
		codec:    codec.Default(),
		errorHandler: func(ctx *Context, err error) {
			log.Printf("[%s] error: %v\n", ctx.Conn().UUID(), err)
		},
	}
}
