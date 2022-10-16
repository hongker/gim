package framework

import (
	"log"
	"sync"
)

// Router
type Router struct {
	rwm      sync.RWMutex
	handlers map[int32]Handler
	codec    Codec
}

// Route register handler for operate
func (router *Router) Route(operate int32, handler Handler) *Router {
	router.rwm.Lock()
	router.handlers[operate] = handler
	router.rwm.Unlock()
	return router
}

func (router *Router) Request(ctx *Context) {
	// unpack
	packet, err := router.codec.Unpack(ctx.body)
	if err != nil {
		log.Println("unpack:", err)
		return
	}
	router.rwm.RLock()
	defer router.rwm.RUnlock()

	// find handler
	handler, ok := router.handlers[packet.Operate]
	if !ok {
		return
	}

	ctx.packet = packet
	response, err := handler(ctx)
	if err != nil {
		return
	}

	packet.Operate++
	// pack response
	msg, err := router.codec.Pack(packet, response)
	if err != nil {
		return
	}
	ctx.Conn().Push(msg)
}

func NewRouter() *Router {
	return &Router{
		handlers: map[int32]Handler{},
		codec:    &DefaultCodec{},
	}
}
