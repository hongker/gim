package gateway

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/socket"
	"github.com/ebar-go/ego/utils/runtime"
	"time"
)

type Callback struct {
	codec    Codec
	em       *EventManager
	provider ProtoProvider
}

func NewCallback() *Callback {
	c := &Callback{
		codec:    DefaultCodec(),
		em:       NewEventManager(),
		provider: NewSharedProtoProvider(),
	}
	return c
}

func (c *Callback) OnConnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Connected, IP: %s", conn.ID(), conn.IP())

	// close the connection if client don't send heartbeat request.
	timer := time.NewTimer(time.Minute)
	go func() {
		<-timer.C
		runtime.HandlerError(conn.Close(), func(err error) {
			component.Provider().Logger().Errorf("[%s] closed failed: %v", conn.ID(), err)
		})
	}()

	conn.Property().Set("timer", timer)
}
func (c *Callback) OnDisconnect(conn socket.Connection) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.ID())
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
