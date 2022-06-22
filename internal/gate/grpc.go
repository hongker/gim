package gate

import (
	"context"
	"gim/api/server"
	"github.com/pkg/errors"
)

func (srv *Server) PushMsg(ctx context.Context, request *server.PushMsgRequest) (res *server.PushMsgResponse, err error) {
	packet, err := request.Proto.Pack()
	if err != nil {
		return
	}
	for _, key := range request.Keys {
		ch := srv.bucket.GetChannel(key)
		if ch == nil {
			continue
		}
		ch.conn.Push(packet)
	}
	return
}

func (srv *Server) BroadcastRoom(ctx context.Context, request *server.BroadcastRoomRequest) (res *server.BroadcastRoomResponse, err error) {
	room := srv.bucket.GetRoom(request.RoomID)
	if room == nil {
		err = errors.New("not exist")
		return
	}
	packet, err := request.Proto.Pack()
	if err != nil {
		return
	}

	for _, ch := range room.channels {
		ch.conn.Push(packet)
	}
	return
}
