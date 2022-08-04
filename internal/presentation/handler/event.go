package handler

import (
	"gim/api"
	"gim/internal/application"
	"gim/internal/domain/dto"
	"gim/internal/domain/event"
	"gim/pkg/network"
	"time"
)

type EventHandler struct {
	expired time.Duration
	gateApp *application.GateApp
}

func (h *EventHandler) RegisterEvents()  {
	event.Listen(event.Connect, h.Connect)
	event.Listen(event.Login, h.Login)
	event.Listen(event.Disconnect, h.Disconnect)
	event.Listen(event.JoinGroup, h.JoinGroup)
	event.Listen(event.LeaveGroup, h.LeaveGroup)
	event.Listen(event.Push, h.Push)
}

func (h *EventHandler) Push(params ...interface{})   {
	if len(params) <= 1 {
		return
	}

	sessionType := params[0].(string)
	targetId := params[1].(string)
	batchMessages := params[2].(*dto.BatchMessage)
	packet := api.BuildPacket(api.OperateMessagePush, batchMessages)
	h.gateApp.Push(sessionType, targetId, packet.Encode())
}

func (h *EventHandler) Connect(params ...interface{}) {
	if len(params) <= 1 {
		return
	}
	conn := params[0].(*network.Connection)
	// 如果用户未按时登录，通过定时任务关闭连接，释放资源
	time.AfterFunc(h.expired, func() {
		if !h.gateApp.CheckConnExist(conn) {
			conn.Close()
		}

	})
}

func (h *EventHandler) Login(params ...interface{}) {
	if len(params) <= 1 {
		return
	}
	uid := params[0].(string)
	conn := params[1].(*network.Connection)
	h.gateApp.RegisterConn(uid, conn)
}
func (h *EventHandler) Disconnect(params ...interface{}) {
	if len(params) < 1 {
		return
	}
	conn := params[0].(*network.Connection)
	h.gateApp.RemoveConn(conn)
}
func (h *EventHandler) JoinGroup(params ...interface{}) {
	if len(params) <= 1 {
		return
	}
	roomId := params[0].(string)
	conn := params[1].(*network.Connection)
	h.gateApp.JoinRoom(roomId, conn)

}
func (h *EventHandler) LeaveGroup(params ...interface{}) {
	if len(params) <= 1 {
		return
	}
	roomId := params[0].(string)
	conn := params[1].(*network.Connection)
	h.gateApp.LeaveRoom(roomId, conn)
}

func NewEventHandler(gateApp *application.GateApp) *EventHandler {
	h :=  &EventHandler{gateApp: gateApp, expired: time.Minute}
	return h
}
