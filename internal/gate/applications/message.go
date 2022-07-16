package applications

import (
	"context"
	"gim/api"
	"gim/api/client"
	"gim/api/protocol"
	"gim/api/server"
	"gim/internal/gate/domain/entity"
	"gim/pkg/network"
	"github.com/pkg/errors"
	"log"
)

type MessageApp interface {
	Push(key string, packet []byte) error
	PushRoom(key string, packet []byte) error
	Auth(ctx context.Context, proto *protocol.Proto, conn *network.Connection) (uid string, err error)
	Receive(ctx context.Context, proto *protocol.Proto, conn *network.Connection) (res *server.ReceiveResponse, err error)
}

type messageApp struct {
	logic server.LogicClient
	bucket *entity.Bucket
}

func newMessageApp() MessageApp {
	return &messageApp{bucket: entity.NewBucket()}
}

func (app messageApp) Push(key string, packet []byte) error {
	ch := app.bucket.GetChannel(key)
	if ch == nil {
		return errors.New("channel not found")
	}
	ch.Conn().Push(packet)
	return nil
}

func (app messageApp) PushRoom(key string, packet []byte) (err error) {
	room := app.bucket.GetRoom(key)
	if room == nil {
		err = errors.New("not exist")
		return
	}

	for _, ch := range room.Channels() {
		ch.Conn().Push(packet)
	}
	return
}

func (app messageApp) Auth(ctx context.Context, proto *protocol.Proto, conn *network.Connection) (uid string, err error) {
	if proto.Op != api.OperateAuth {
		err = errors.New("invalid operate")
		return
	}

	req := new(client.AuthRequest)
	if err = proto.Bind(req); err != nil {
		err = errors.WithMessage(err, "invalid request")
		return
	}

	// 验证用户
	res, err := app.logic.Auth(ctx, &server.AuthRequest{
		AppId: req.AppId,
		Name:  req.Name,
	})
	if err != nil {
		return
	}

	// 验证通过后，将channel加入到bucket,用于服务端广播消息
	app.bucket.AddChannel(entity.NewChannel(res.Uid, conn))
	uid = res.Uid
	return
}


func (app messageApp) Receive(ctx context.Context, proto *protocol.Proto, conn *network.Connection) (res *server.ReceiveResponse, err error) {
	channel := app.bucket.GetChannel(conn.ID())

	switch proto.Op {
	case api.OperateGroupJoin: // 加群
		groupJoinRequest := new(client.GroupJoinRequest)
		if err = proto.Bind(groupJoinRequest); err != nil {
			err = errors.WithMessage(err, "invalid request")
			return
		}
		app.bucket.GetRoom(groupJoinRequest.GroupId).Add(channel)
	case api.OperateGroupQuit: // 退群
		groupQuitRequest := new(client.GroupQuitRequest)
		if err = proto.Bind(groupQuitRequest); err != nil {
			err = errors.WithMessage(err, "invalid request")
			return
		}
		app.bucket.GetRoom(groupQuitRequest.GroupId).Remove(channel)
	}

	log.Println("receive:", proto.Op)
	return app.logic.Receive(ctx, &server.ReceiveRequest{
		Mid:   channel.Key(),
		Proto: proto,
	})
}







