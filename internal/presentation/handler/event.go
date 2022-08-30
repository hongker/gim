package handler

import (
	"gim/api"
	"gim/internal/domain/event"
	"gim/internal/domain/types"
	"time"
)

type EventHandler struct {
	expired    time.Duration
	collection *types.Collection
}

func (h *EventHandler) Push(ev *event.PushMessageEvent) {
	packet := api.BuildPacket(api.OperateMessagePush, ev.BatchMessage)
	h.collection.Push(ev.SessionType, ev.TargetId, packet.Encode())
}

func (h *EventHandler) Connect(ev *event.ConnectEvent) {
	// 如果用户未按时登录，通过定时任务关闭连接，释放资源
	h.collection.Add(ev.Connection)
	h.collection.Refresh(ev.Connection, h.expired)
}

func (h *EventHandler) Heartbeat(ev *event.HeartbeatEvent) {
	h.collection.Refresh(ev.Connection, h.expired)
}

func (h *EventHandler) Login(ev *event.LoginEvent) {
	h.collection.RegisterConn(ev.UserId, ev.Connection)
}
func (h *EventHandler) Disconnect(ev *event.DisconnectEvent) {
	h.collection.RemoveConn(ev.Connection)
}
func (h *EventHandler) JoinGroup(ev *event.JoinGroupEvent) {

	h.collection.JoinRoom(ev.GroupId, ev.Connection)

}
func (h *EventHandler) LeaveGroup(ev *event.LeaveGroupEvent) {

	h.collection.LeaveRoom(ev.GroupId, ev.Connection)
}

func NewEventHandler(collection *types.Collection, expired time.Duration) *EventHandler {
	return &EventHandler{collection: collection, expired: expired}
}
