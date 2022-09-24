package socket

import (
	"gim/internal/module/gateway/application"
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/server/ws"
	"github.com/ebar-go/ego/utils/runtime"
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
	defer runtime.HandleCrash()
	component.Provider().Logger().Infof("OnMessage: %s", string(ctx.Body()))

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

func (c *Callback) matchEvents(proto *Proto) Event {
	return c.events[proto.OperateType()]
}

func (c *Callback) prepare() {
	c.initHandler()
}

func (c *Callback) initHandler() {
	em := &EventManager{
		userApp: application.NewUserApplication(),
	}
	c.events[LoginOperate] = Action[dto.UserLoginRequest, dto.UserLoginResponse](em.Login)
	c.events[HeartbeatOperate] = Action[dto.SocketHeartbeatRequest, dto.SocketHeartbeatResponse](em.Heartbeat)
}
