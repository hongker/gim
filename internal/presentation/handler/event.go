package handler

import (
	"gim/api"
	"gim/internal/domain/dto"
	"gim/internal/domain/types"
	"gim/pkg/network"
	"time"
)

type EventHandler struct {
	expired    time.Duration
	collection *types.Collection
}

func (h *EventHandler) Push(params ...interface{}) {
	if len(params) < 3 {
		return
	}

	sessionType := params[0].(string)
	targetId := params[1].(string)
	batchMessages := params[2].(*dto.BatchMessage)
	packet := api.BuildPacket(api.OperateMessagePush, batchMessages)
	h.collection.Push(sessionType, targetId, packet.Encode())
}

func (h *EventHandler) Connect(params ...interface{}) {
	if len(params) < 1 {
		return
	}
	conn := params[0].(*network.Connection)
	// 如果用户未按时登录，通过定时任务关闭连接，释放资源
	h.collection.Add(conn)
	h.collection.Refresh(conn, h.expired)
}

func (h *EventHandler) Heartbeat(params ...interface{}) {
	if len(params) < 1 {
		return
	}
	conn := params[0].(*network.Connection)
	h.collection.Refresh(conn, h.expired)
}

func (h *EventHandler) Login(params ...interface{}) {
	if len(params) < 2 {
		return
	}
	uid := params[0].(string)
	conn := params[1].(*network.Connection)
	h.collection.RegisterConn(uid, conn)
}
func (h *EventHandler) Disconnect(params ...interface{}) {
	if len(params) < 1 {
		return
	}
	conn := params[0].(*network.Connection)
	h.collection.RemoveConn(conn)
}
func (h *EventHandler) JoinGroup(params ...interface{}) {
	if len(params) < 2 {
		return
	}
	roomId := params[0].(string)
	conn := params[1].(*network.Connection)
	h.collection.JoinRoom(roomId, conn)

}
func (h *EventHandler) LeaveGroup(params ...interface{}) {
	if len(params) < 2 {
		return
	}
	roomId := params[0].(string)
	conn := params[1].(*network.Connection)
	h.collection.LeaveRoom(roomId, conn)
}

func NewEventHandler(collection *types.Collection, expired time.Duration) *EventHandler {
	return &EventHandler{collection: collection, expired: expired}
}
