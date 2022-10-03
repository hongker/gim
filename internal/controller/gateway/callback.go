package gateway

import (
	"gim/api"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/socket"
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

type Callback struct {
	codec    api.Codec
	em       *EventManager
	provider api.ProtoProvider
}

func NewCallback(heartbeatInterval time.Duration) *Callback {
	c := &Callback{
		codec:    api.DefaultCodec(),
		em:       NewEventManager(heartbeatInterval),
		provider: api.NewSharedProtoProvider(),
	}
	return c
}

func (c *Callback) OnConnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Connected, IP: %s", conn.ID(), conn.IP())
	c.em.InitializeConn(conn)

}
func (c *Callback) OnDisconnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.ID())
	c.em.FinalizeConn(conn)
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

	c.em.Handle(ctx, proto)

	response := c.codec.Encode(proto)
	component.Provider().Logger().Infof("[%s] OnResponse: %s", ctx.Conn().ID(), string(response))
	ctx.Output(response)
}
