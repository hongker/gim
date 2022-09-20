package socket

import (
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/component"
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

func (c *Callback) OnConnect(conn ws.Conn) {
	component.Provider().Logger().Infof("Connect: %s", conn.IP())
}
func (c *Callback) OnDisconnect(conn ws.Conn) {
	component.Provider().Logger().Infof("Disconnect: %s", conn.IP())
}
func (c *Callback) OnMessage(ctx *ws.Context) {
	component.Provider().Logger().Infof("OnMessage: %s", string(ctx.Body()))
	defer c.handleCrash(ctx)

	proto, err := c.codec.Decode(ctx.Body())
	if err != nil {
		return
	}

	handler := c.matchEvents(proto)
	if handler == nil {
		return
	}

	handler(ctx, proto)
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
	c.events[LoginOperate] = Action[dto.SocketLoginRequest, dto.SocketLoginResponse](LoginEvent)
	c.events[HeartbeatOperate] = Action[dto.SocketHeartbeatRequest, dto.SocketHeartbeatResponse](HeartbeatEvent)
}
