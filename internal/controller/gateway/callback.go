package gateway

import (
	"gim/api"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/socket"
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

type Callback struct {
	// codec decode request and encode response
	codec api.Codec

	// handler handle operation
	handler  *Handler
	provider api.ProtoProvider
}

func NewCallback(heartbeatInterval time.Duration) *Callback {
	return &Callback{
		codec:    api.DefaultCodec(),
		handler:  NewHandler(heartbeatInterval),
		provider: api.NewSharedProtoProvider(),
	}
}

func (c *Callback) OnConnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Connected, IP: %s", conn.ID(), conn.IP())

	// initialize connection
	c.handler.InitializeConn(conn)

}
func (c *Callback) OnDisconnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.ID())

	// finalize connection
	c.handler.FinalizeConn(conn)
}

func (c *Callback) OnMessage(ctx *socket.Context) {
	defer runtime.HandleCrash()
	component.Provider().Logger().Infof("[%s] OnMessage: %s", ctx.Conn().ID(), string(ctx.Body()))

	// acquire proto from provider,optimize for GC.
	proto := c.provider.Acquire()
	// release proto to provider
	defer c.provider.Release(proto)

	err := c.codec.Decode(ctx.Body(), proto)
	if err != nil {
		component.Provider().Logger().Errorf("[%s] OnDecode: %v", ctx.Conn().ID(), err)
		return
	}

	c.handler.Handle(ctx, proto)

	response := c.codec.Encode(proto)
	component.Provider().Logger().Infof("[%s] OnResponse: %s", ctx.Conn().ID(), string(response))
	ctx.Output(response)
}
