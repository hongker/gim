package framework

import "sync"

type Router struct {
	rwm      sync.RWMutex
	handlers map[int]Handler
	codec    Codec
}

func NewRouter() *Router {
	return &Router{
		handlers: map[int]Handler{},
	}
}

func (router *Router) Handle(operate int, handler Handler) *Router {
	router.rwm.Lock()
	router.handlers[operate] = handler
	router.rwm.Unlock()
	return router
}

func (router *Router) Request() HandleFunc {
	return func(ctx *Context) {
		operate, err := router.codec.Unpack(ctx.body)
		if err != nil {
			return
		}
		router.rwm.RLock()
		defer router.rwm.RUnlock()

		handler, ok := router.handlers[operate]
		if !ok {
			return
		}
		response, err := handler(ctx, router.codec.Serializer())
		if err != nil {
			return
		}

		msg, err := router.codec.Pack(operate+1, response)
		if err != nil {
			return
		}
		ctx.Conn().Push(msg)
	}
}
