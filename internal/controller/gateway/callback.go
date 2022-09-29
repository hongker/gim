package gateway

import (
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/ws"
	"github.com/ebar-go/ego/utils/runtime"
)

type Callback struct {
	codec Codec
	em    *EventManager
}

func NewCallback() *Callback {
	c := &Callback{
		codec: DefaultCodec(),
		em:    NewEventManager(),
	}
	return c
}

func (c *Callback) OnConnect(conn ws.Conn) {
	component.Provider().Logger().Infof("[%s] Connected, IP: %s", conn.ID(), conn.IP())
}
func (c *Callback) OnDisconnect(conn ws.Conn) {
	component.Provider().Logger().Infof("[%s] Disconnected", conn.ID())
}
func (c *Callback) OnMessage(ctx *ws.Context) {
	defer runtime.HandleCrash()
	component.Provider().Logger().Infof("[%s] OnMessage: %s", ctx.Conn().ID(), string(ctx.Body()))

	proto, err := c.codec.Decode(ctx.Body())
	if err != nil {
		component.Provider().Logger().Errorf("[%s] OnDecode: %v", ctx.Conn().ID(), err)
		return
	}

	c.em.Handle(ctx, proto)

	response := c.codec.Encode(proto)
	component.Provider().Logger().Infof("[%s] OnResponse: %s", ctx.Conn().ID(), string(response))
	ctx.Output(response)
}
