package socket

import (
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/ego/server/ws"
)

type Callback struct {
	codec  Codec
	events map[OperateType]Event
}

func NewCallback() *Callback {
	c := &Callback{
		codec:  DefaultCodec(),
		events: map[OperateType]Event{},
	}
	c.prepare()
	return c
}

func (c *Callback) OnConnect(conn ws.Conn)    {}
func (c *Callback) OnDisconnect(conn ws.Conn) {}
func (c *Callback) OnMessage(ctx *ws.Context) {
	defer c.handleCrash(ctx)

	proto, err := c.codec.Decode(ctx.Body())
	if err != nil {
		return
	}

	handler := c.matchEvents(proto)
	if handler == nil {
		return
	}

	err = handler(ctx, proto)
	if err != nil {
		return
	}
	ctx.Output(c.codec.Encode(proto))
}

func (c *Callback) handleCrash(ctx *ws.Context) {
	if err := recover(); err != nil {
		switch err.(type) {
		case errors.Error:
		default:

		}
	}
}

func (c *Callback) matchEvents(proto *Proto) Event {
	return c.events[proto.OperateType()]
}

func (c *Callback) prepare() {
	c.initHandler()
}

func (c *Callback) initHandler() {
	c.events[ConnectOperate] = Action[dto.ConnectRequest, dto.ConnectResponse](ConnectEvent)
	c.events[DisconnectOperate] = Action[dto.DisconnectRequest, dto.DisconnectResponse](DisconnectEvent)
	c.events[HeartbeatOperate] = Action[dto.HeartbeatRequest, dto.HeartbeatResponse](HeartbeatEvent)
}
