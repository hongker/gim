package framework

type Router struct {
	handlers map[int]Handler
	codec    Codec
}

func NewRouter() *Router { return &Router{} }

func (router *Router) Handle(operate int, handler Handler) *Router { return router }
func (router *Router) Request() HandleFunc {
	return func(ctx *Context) {
		operate := router.codec.Unpack(ctx.body)
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
