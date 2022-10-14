package framework

type Router struct{}

func NewRouter() *Router { return &Router{} }

func (router *Router) Handle(operate int, handler Handler) *Router { return router }
